package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/jung-kurt/gofpdf"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// ItineraryDeliveryService handles PDF generation, ICS files, and delivery
type ItineraryDeliveryService struct {
	emailConfig   *EmailConfig
	smsConfig     *SMSConfig
	storageConfig *StorageConfig
	templateDir   string
	firebase      *FirebaseService
}

// EmailConfig contains email service configuration
type EmailConfig struct {
	SMTPHost  string
	SMTPPort  int
	Username  string
	Password  string
	FromEmail string
	FromName  string
	Enabled   bool
}

// SMSConfig contains SMS service configuration
type SMSConfig struct {
	TwilioAccountSID  string
	TwilioAuthToken   string
	TwilioPhoneNumber string
	Enabled           bool
}

// StorageConfig contains file storage configuration
type StorageConfig struct {
	BasePath     string
	BaseURL      string
	CloudStorage bool
	BucketName   string
}

// NewItineraryDeliveryService creates a new delivery service
func NewItineraryDeliveryService(
	emailConfig *EmailConfig,
	smsConfig *SMSConfig,
	storageConfig *StorageConfig,
	firebase *FirebaseService,
) *ItineraryDeliveryService {
	return &ItineraryDeliveryService{
		emailConfig:   emailConfig,
		smsConfig:     smsConfig,
		storageConfig: storageConfig,
		templateDir:   "templates",
		firebase:      firebase,
	}
}

// DeliveryFormat represents the format for itinerary delivery
type DeliveryFormat string

const (
	FormatPDF  DeliveryFormat = "pdf"
	FormatICS  DeliveryFormat = "ics"
	FormatJSON DeliveryFormat = "json"
	FormatHTML DeliveryFormat = "html"
)

// DeliveryMethod represents how the itinerary should be delivered
type DeliveryMethod string

const (
	MethodEmail    DeliveryMethod = "email"
	MethodSMS      DeliveryMethod = "sms"
	MethodDownload DeliveryMethod = "download"
	MethodPush     DeliveryMethod = "push"
)

// DeliveryRequest represents a request to deliver an itinerary
type DeliveryRequest struct {
	TripID          string         `json:"trip_id"`
	UserID          string         `json:"user_id"`
	Format          DeliveryFormat `json:"format"`
	Method          DeliveryMethod `json:"method"`
	Recipient       string         `json:"recipient"` // email or phone number
	Language        string         `json:"language"`
	IncludeBookings bool           `json:"include_bookings"`
	IncludeMap      bool           `json:"include_map"`
	CustomMessage   string         `json:"custom_message,omitempty"`
	Template        string         `json:"template,omitempty"`
}

// DeliveryResult represents the result of a delivery operation
type DeliveryResult struct {
	DeliveryID    string     `json:"delivery_id"`
	TripID        string     `json:"trip_id"`
	UserID        string     `json:"user_id"`
	Format        string     `json:"format"`
	Method        string     `json:"method"`
	FileURL       string     `json:"file_url,omitempty"`
	FileName      string     `json:"file_name,omitempty"`
	Status        string     `json:"status"` // success, failed, pending
	DeliveredAt   time.Time  `json:"delivered_at"`
	ErrorMessage  string     `json:"error_message,omitempty"`
	DownloadCount int        `json:"download_count"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
}

// ItineraryData represents structured itinerary data for delivery
type ItineraryData struct {
	TripID            string               `json:"trip_id"`
	Destination       string               `json:"destination"`
	StartDate         time.Time            `json:"start_date"`
	EndDate           time.Time            `json:"end_date"`
	Travelers         int                  `json:"travelers"`
	Budget            float64              `json:"budget"`
	Currency          string               `json:"currency"`
	Title             string               `json:"title"`
	Description       string               `json:"description"`
	DailyItinerary    map[int]DayItinerary `json:"daily_itinerary"`
	Hotels            []HotelBooking       `json:"hotels,omitempty"`
	Transportation    []TransportBooking   `json:"transportation,omitempty"`
	Activities        []ActivityBooking    `json:"activities,omitempty"`
	ImportantInfo     []string             `json:"important_info"`
	EmergencyContacts []EmergencyContact   `json:"emergency_contacts"`
	TotalCost         float64              `json:"total_cost"`
	CreatedAt         time.Time            `json:"created_at"`
	LastModified      time.Time            `json:"last_modified"`
}

// DayItinerary represents a single day's activities
type DayItinerary struct {
	Date      time.Time    `json:"date"`
	DayNumber int          `json:"day_number"`
	Title     string       `json:"title"`
	Morning   []Activity   `json:"morning"`
	Afternoon []Activity   `json:"afternoon"`
	Evening   []Activity   `json:"evening"`
	Meals     []Meal       `json:"meals"`
	Notes     string       `json:"notes,omitempty"`
	Weather   *WeatherInfo `json:"weather,omitempty"`
	TotalCost float64      `json:"total_cost"`
}

// Activity represents a single activity
type Activity struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Location    Location  `json:"location"`
	Description string    `json:"description"`
	Cost        float64   `json:"cost"`
	BookingRef  string    `json:"booking_ref,omitempty"`
	Status      string    `json:"status"` // confirmed, pending, optional
	Tips        []string  `json:"tips,omitempty"`
}

// Meal represents a meal/dining activity
type Meal struct {
	Type       string    `json:"type"` // breakfast, lunch, dinner, snack
	Restaurant string    `json:"restaurant"`
	Location   Location  `json:"location"`
	Time       time.Time `json:"time"`
	Cost       float64   `json:"cost"`
	Cuisine    string    `json:"cuisine"`
	BookingRef string    `json:"booking_ref,omitempty"`
}

// HotelBooking represents hotel booking information
type HotelBooking struct {
	Name            string    `json:"name"`
	Address         string    `json:"address"`
	CheckIn         time.Time `json:"check_in"`
	CheckOut        time.Time `json:"check_out"`
	RoomType        string    `json:"room_type"`
	Nights          int       `json:"nights"`
	TotalCost       float64   `json:"total_cost"`
	ConfirmationNum string    `json:"confirmation_number"`
	Contact         string    `json:"contact"`
	Amenities       []string  `json:"amenities"`
}

// TransportBooking represents transportation booking
type TransportBooking struct {
	Type          string    `json:"type"` // flight, train, bus, taxi
	From          string    `json:"from"`
	To            string    `json:"to"`
	DepartureTime time.Time `json:"departure_time"`
	ArrivalTime   time.Time `json:"arrival_time"`
	Provider      string    `json:"provider"`
	BookingRef    string    `json:"booking_ref"`
	SeatNumber    string    `json:"seat_number,omitempty"`
	Cost          float64   `json:"cost"`
	Status        string    `json:"status"`
}

// ActivityBooking represents activity booking information
type ActivityBooking struct {
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Date         time.Time `json:"date"`
	Time         time.Time `json:"time"`
	Duration     string    `json:"duration"`
	Participants int       `json:"participants"`
	Cost         float64   `json:"cost"`
	BookingRef   string    `json:"booking_ref"`
	MeetingPoint string    `json:"meeting_point"`
	Instructions []string  `json:"instructions"`
}

// EmergencyContact represents emergency contact information
type EmergencyContact struct {
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
	Phone        string `json:"phone"`
	Email        string `json:"email,omitempty"`
	Available24h bool   `json:"available_24h"`
}

// WeatherInfo represents weather information for a day
type WeatherInfo struct {
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
}

// GenerateAndDeliverItinerary generates and delivers an itinerary
func (d *ItineraryDeliveryService) GenerateAndDeliverItinerary(ctx context.Context, req *DeliveryRequest) (*DeliveryResult, error) {
	// Get itinerary data
	itineraryData, err := d.getItineraryData(ctx, req.TripID, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get itinerary data: %w", err)
	}

	// Generate file based on format
	fileData, fileName, err := d.generateFile(ctx, itineraryData, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate file: %w", err)
	}

	// Store file
	fileURL, err := d.storeFile(ctx, fileData, fileName, req.TripID)
	if err != nil {
		return nil, fmt.Errorf("failed to store file: %w", err)
	}

	// Create delivery result
	result := &DeliveryResult{
		DeliveryID:  d.generateDeliveryID(req.TripID, req.UserID),
		TripID:      req.TripID,
		UserID:      req.UserID,
		Format:      string(req.Format),
		Method:      string(req.Method),
		FileURL:     fileURL,
		FileName:    fileName,
		Status:      "pending",
		DeliveredAt: time.Now(),
	}

	// Deliver based on method
	switch req.Method {
	case MethodEmail:
		err = d.deliverByEmail(ctx, req, fileURL, fileName, itineraryData)
	case MethodSMS:
		err = d.deliverBySMS(ctx, req, fileURL, fileName)
	case MethodDownload:
		// No additional delivery needed
		result.Status = "success"
	case MethodPush:
		err = d.deliverByPush(ctx, req, fileURL, fileName)
	default:
		err = fmt.Errorf("unsupported delivery method: %s", req.Method)
	}

	if err != nil {
		result.Status = "failed"
		result.ErrorMessage = err.Error()
	} else {
		result.Status = "success"
	}

	// Store delivery record
	d.storeDeliveryRecord(ctx, result)

	return result, err
}

// generateFile generates the file in the requested format
func (d *ItineraryDeliveryService) generateFile(ctx context.Context, data *ItineraryData, req *DeliveryRequest) ([]byte, string, error) {
	switch req.Format {
	case FormatPDF:
		return d.generatePDF(data, req)
	case FormatICS:
		return d.generateICS(data, req)
	case FormatJSON:
		return d.generateJSON(data, req)
	case FormatHTML:
		return d.generateHTML(data, req)
	default:
		return nil, "", fmt.Errorf("unsupported format: %s", req.Format)
	}
}

// generatePDF creates a PDF version of the itinerary
func (d *ItineraryDeliveryService) generatePDF(data *ItineraryData, req *DeliveryRequest) ([]byte, string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set fonts
	pdf.SetFont("Arial", "B", 16)

	// Title
	pdf.Cell(0, 10, fmt.Sprintf("Travel Itinerary - %s", data.Destination))
	pdf.Ln(15)

	// Trip details
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Trip Dates: %s to %s",
		data.StartDate.Format("January 2, 2006"),
		data.EndDate.Format("January 2, 2006")))
	pdf.Ln(8)

	pdf.Cell(0, 8, fmt.Sprintf("Travelers: %d", data.Travelers))
	pdf.Ln(8)

	pdf.Cell(0, 8, fmt.Sprintf("Total Budget: %.2f %s", data.Budget, data.Currency))
	pdf.Ln(15)

	// Daily itinerary
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Daily Itinerary")
	pdf.Ln(12)

	for dayNum := 1; dayNum <= len(data.DailyItinerary); dayNum++ {
		if dayData, exists := data.DailyItinerary[dayNum]; exists {
			d.addDayToPDF(pdf, dayNum, dayData)
		}
	}

	// Hotels section
	if len(data.Hotels) > 0 {
		d.addHotelsToPDF(pdf, data.Hotels)
	}

	// Transportation section
	if len(data.Transportation) > 0 {
		d.addTransportationToPDF(pdf, data.Transportation)
	}

	// Important information
	if len(data.ImportantInfo) > 0 {
		d.addImportantInfoToPDF(pdf, data.ImportantInfo)
	}

	// Emergency contacts
	if len(data.EmergencyContacts) > 0 {
		d.addEmergencyContactsToPDF(pdf, data.EmergencyContacts)
	}

	// Generate file data
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	fileName := fmt.Sprintf("itinerary_%s_%s.pdf", data.TripID, time.Now().Format("20060102"))
	return buf.Bytes(), fileName, nil
}

// generateICS creates an ICS calendar file
func (d *ItineraryDeliveryService) generateICS(data *ItineraryData, req *DeliveryRequest) ([]byte, string, error) {
	var ics strings.Builder

	// ICS header
	ics.WriteString("BEGIN:VCALENDAR\r\n")
	ics.WriteString("VERSION:2.0\r\n")
	ics.WriteString("PRODID:-//AuraTravel//AuraTravel AI//EN\r\n")
	ics.WriteString("CALSCALE:GREGORIAN\r\n")
	ics.WriteString("METHOD:PUBLISH\r\n")

	// Add each activity as an event
	for _, dayData := range data.DailyItinerary {
		d.addActivitiesToICS(&ics, dayData.Morning, data.TripID)
		d.addActivitiesToICS(&ics, dayData.Afternoon, data.TripID)
		d.addActivitiesToICS(&ics, dayData.Evening, data.TripID)
	}

	// Add hotel check-ins/check-outs
	for _, hotel := range data.Hotels {
		d.addHotelToICS(&ics, hotel, data.TripID)
	}

	// Add transportation
	for _, transport := range data.Transportation {
		d.addTransportToICS(&ics, transport, data.TripID)
	}

	// ICS footer
	ics.WriteString("END:VCALENDAR\r\n")

	fileName := fmt.Sprintf("itinerary_%s_%s.ics", data.TripID, time.Now().Format("20060102"))
	return []byte(ics.String()), fileName, nil
}

// generateJSON creates a JSON version of the itinerary
func (d *ItineraryDeliveryService) generateJSON(data *ItineraryData, req *DeliveryRequest) ([]byte, string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fileName := fmt.Sprintf("itinerary_%s_%s.json", data.TripID, time.Now().Format("20060102"))
	return jsonData, fileName, nil
}

// generateHTML creates an HTML version of the itinerary
func (d *ItineraryDeliveryService) generateHTML(data *ItineraryData, req *DeliveryRequest) ([]byte, string, error) {
	// Load template
	templatePath := filepath.Join(d.templateDir, "itinerary.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		// Fallback to inline template
		tmpl, err = template.New("itinerary").Parse(d.getDefaultHTMLTemplate())
		if err != nil {
			return nil, "", fmt.Errorf("failed to parse template: %w", err)
		}
	}

	// Execute template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, "", fmt.Errorf("failed to execute template: %w", err)
	}

	fileName := fmt.Sprintf("itinerary_%s_%s.html", data.TripID, time.Now().Format("20060102"))
	return buf.Bytes(), fileName, nil
}

// deliverByEmail sends the itinerary via email
func (d *ItineraryDeliveryService) deliverByEmail(ctx context.Context, req *DeliveryRequest, fileURL, fileName string, data *ItineraryData) error {
	if !d.emailConfig.Enabled {
		return fmt.Errorf("email delivery not enabled")
	}

	// Get user email if not provided
	recipient := req.Recipient
	if recipient == "" {
		userEmail, err := d.getUserEmail(ctx, req.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user email: %w", err)
		}
		recipient = userEmail
	}

	// Prepare email content
	subject := fmt.Sprintf("Your Travel Itinerary - %s", data.Destination)
	body := d.buildEmailBody(data, fileURL, req.CustomMessage)

	// Send email
	return d.sendEmail(recipient, subject, body, fileURL, fileName)
}

// deliverBySMS sends a download link via SMS
func (d *ItineraryDeliveryService) deliverBySMS(ctx context.Context, req *DeliveryRequest, fileURL, fileName string) error {
	if !d.smsConfig.Enabled {
		return fmt.Errorf("SMS delivery not enabled")
	}

	// Get user phone if not provided
	recipient := req.Recipient
	if recipient == "" {
		userPhone, err := d.getUserPhone(ctx, req.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user phone: %w", err)
		}
		recipient = userPhone
	}

	// Prepare SMS content
	message := fmt.Sprintf("Your travel itinerary is ready! Download it here: %s", fileURL)
	if req.CustomMessage != "" {
		message = req.CustomMessage + "\n\n" + message
	}

	// Send SMS
	return d.sendSMS(recipient, message)
}

// deliverByPush sends a push notification with download link
func (d *ItineraryDeliveryService) deliverByPush(ctx context.Context, req *DeliveryRequest, fileURL, fileName string) error {
	// This would integrate with the NotificationService
	// For now, return success
	log.Printf("Push notification delivery for trip %s: %s", req.TripID, fileURL)
	return nil
}

// Helper methods for file generation

func (d *ItineraryDeliveryService) addDayToPDF(pdf *gofpdf.Fpdf, dayNum int, dayData DayItinerary) {
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Day %d - %s", dayNum, dayData.Date.Format("Monday, January 2")))
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)

	if len(dayData.Morning) > 0 {
		pdf.Cell(40, 6, "Morning:")
		pdf.Ln(6)
		for _, activity := range dayData.Morning {
			pdf.Cell(10, 5, "")
			pdf.Cell(0, 5, fmt.Sprintf("• %s (%s)", activity.Name, activity.StartTime.Format("3:04 PM")))
			pdf.Ln(5)
		}
		pdf.Ln(3)
	}

	if len(dayData.Afternoon) > 0 {
		pdf.Cell(40, 6, "Afternoon:")
		pdf.Ln(6)
		for _, activity := range dayData.Afternoon {
			pdf.Cell(10, 5, "")
			pdf.Cell(0, 5, fmt.Sprintf("• %s (%s)", activity.Name, activity.StartTime.Format("3:04 PM")))
			pdf.Ln(5)
		}
		pdf.Ln(3)
	}

	if len(dayData.Evening) > 0 {
		pdf.Cell(40, 6, "Evening:")
		pdf.Ln(6)
		for _, activity := range dayData.Evening {
			pdf.Cell(10, 5, "")
			pdf.Cell(0, 5, fmt.Sprintf("• %s (%s)", activity.Name, activity.StartTime.Format("3:04 PM")))
			pdf.Ln(5)
		}
		pdf.Ln(3)
	}

	pdf.Ln(5)
}

func (d *ItineraryDeliveryService) addHotelsToPDF(pdf *gofpdf.Fpdf, hotels []HotelBooking) {
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Accommodation")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)
	for _, hotel := range hotels {
		pdf.Cell(0, 6, fmt.Sprintf("Hotel: %s", hotel.Name))
		pdf.Ln(6)
		pdf.Cell(0, 6, fmt.Sprintf("Address: %s", hotel.Address))
		pdf.Ln(6)
		pdf.Cell(0, 6, fmt.Sprintf("Check-in: %s | Check-out: %s",
			hotel.CheckIn.Format("Jan 2, 2006"),
			hotel.CheckOut.Format("Jan 2, 2006")))
		pdf.Ln(6)
		if hotel.ConfirmationNum != "" {
			pdf.Cell(0, 6, fmt.Sprintf("Confirmation: %s", hotel.ConfirmationNum))
			pdf.Ln(6)
		}
		pdf.Ln(3)
	}
}

func (d *ItineraryDeliveryService) addTransportationToPDF(pdf *gofpdf.Fpdf, transportation []TransportBooking) {
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Transportation")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)
	for _, transport := range transportation {
		pdf.Cell(0, 6, fmt.Sprintf("%s: %s to %s",
			strings.Title(transport.Type), transport.From, transport.To))
		pdf.Ln(6)
		pdf.Cell(0, 6, fmt.Sprintf("Departure: %s | Arrival: %s",
			transport.DepartureTime.Format("Jan 2, 3:04 PM"),
			transport.ArrivalTime.Format("Jan 2, 3:04 PM")))
		pdf.Ln(6)
		if transport.BookingRef != "" {
			pdf.Cell(0, 6, fmt.Sprintf("Booking Reference: %s", transport.BookingRef))
			pdf.Ln(6)
		}
		pdf.Ln(3)
	}
}

func (d *ItineraryDeliveryService) addImportantInfoToPDF(pdf *gofpdf.Fpdf, info []string) {
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Important Information")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)
	for _, item := range info {
		pdf.Cell(0, 6, fmt.Sprintf("• %s", item))
		pdf.Ln(6)
	}
	pdf.Ln(5)
}

func (d *ItineraryDeliveryService) addEmergencyContactsToPDF(pdf *gofpdf.Fpdf, contacts []EmergencyContact) {
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Emergency Contacts")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)
	for _, contact := range contacts {
		pdf.Cell(0, 6, fmt.Sprintf("%s (%s): %s", contact.Name, contact.Relationship, contact.Phone))
		pdf.Ln(6)
	}
}

func (d *ItineraryDeliveryService) addActivitiesToICS(ics *strings.Builder, activities []Activity, tripID string) {
	for i, activity := range activities {
		ics.WriteString("BEGIN:VEVENT\r\n")
		ics.WriteString(fmt.Sprintf("UID:%s-%d@auratravel.com\r\n", tripID, i))
		ics.WriteString(fmt.Sprintf("DTSTART:%s\r\n", activity.StartTime.UTC().Format("20060102T150405Z")))
		ics.WriteString(fmt.Sprintf("DTEND:%s\r\n", activity.EndTime.UTC().Format("20060102T150405Z")))
		ics.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", activity.Name))
		ics.WriteString(fmt.Sprintf("DESCRIPTION:%s\r\n", activity.Description))
		ics.WriteString(fmt.Sprintf("LOCATION:%s\r\n", activity.Location.Address))
		ics.WriteString("END:VEVENT\r\n")
	}
}

func (d *ItineraryDeliveryService) addHotelToICS(ics *strings.Builder, hotel HotelBooking, tripID string) {
	// Check-in event
	ics.WriteString("BEGIN:VEVENT\r\n")
	ics.WriteString(fmt.Sprintf("UID:%s-checkin@auratravel.com\r\n", tripID))
	ics.WriteString(fmt.Sprintf("DTSTART:%s\r\n", hotel.CheckIn.UTC().Format("20060102T150405Z")))
	ics.WriteString(fmt.Sprintf("SUMMARY:Hotel Check-in - %s\r\n", hotel.Name))
	ics.WriteString(fmt.Sprintf("LOCATION:%s\r\n", hotel.Address))
	ics.WriteString("END:VEVENT\r\n")

	// Check-out event
	ics.WriteString("BEGIN:VEVENT\r\n")
	ics.WriteString(fmt.Sprintf("UID:%s-checkout@auratravel.com\r\n", tripID))
	ics.WriteString(fmt.Sprintf("DTSTART:%s\r\n", hotel.CheckOut.UTC().Format("20060102T150405Z")))
	ics.WriteString(fmt.Sprintf("SUMMARY:Hotel Check-out - %s\r\n", hotel.Name))
	ics.WriteString(fmt.Sprintf("LOCATION:%s\r\n", hotel.Address))
	ics.WriteString("END:VEVENT\r\n")
}

func (d *ItineraryDeliveryService) addTransportToICS(ics *strings.Builder, transport TransportBooking, tripID string) {
	ics.WriteString("BEGIN:VEVENT\r\n")
	ics.WriteString(fmt.Sprintf("UID:%s-transport@auratravel.com\r\n", tripID))
	ics.WriteString(fmt.Sprintf("DTSTART:%s\r\n", transport.DepartureTime.UTC().Format("20060102T150405Z")))
	ics.WriteString(fmt.Sprintf("DTEND:%s\r\n", transport.ArrivalTime.UTC().Format("20060102T150405Z")))
	ics.WriteString(fmt.Sprintf("SUMMARY:%s - %s to %s\r\n", strings.Title(transport.Type), transport.From, transport.To))
	ics.WriteString(fmt.Sprintf("DESCRIPTION:Provider: %s\\nBooking: %s\r\n", transport.Provider, transport.BookingRef))
	ics.WriteString("END:VEVENT\r\n")
}

// Email and SMS sending methods

func (d *ItineraryDeliveryService) sendEmail(to, subject, body, attachmentURL, attachmentName string) error {
	// Set up authentication
	auth := smtp.PlainAuth("", d.emailConfig.Username, d.emailConfig.Password, d.emailConfig.SMTPHost)

	// Build message
	msg := d.buildEmailMessage(to, subject, body, attachmentURL, attachmentName)

	// Send email
	addr := fmt.Sprintf("%s:%d", d.emailConfig.SMTPHost, d.emailConfig.SMTPPort)
	return smtp.SendMail(addr, auth, d.emailConfig.FromEmail, []string{to}, []byte(msg))
}

func (d *ItineraryDeliveryService) buildEmailMessage(to, subject, body, attachmentURL, attachmentName string) string {
	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", d.emailConfig.FromName, d.emailConfig.FromEmail))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	if attachmentURL != "" {
		msg.WriteString(fmt.Sprintf("\r\n\r\nDownload your itinerary: <a href=\"%s\">%s</a>", attachmentURL, attachmentName))
	}

	return msg.String()
}

func (d *ItineraryDeliveryService) sendSMS(to, message string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: d.smsConfig.TwilioAccountSID,
		Password: d.smsConfig.TwilioAuthToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetFrom(d.smsConfig.TwilioPhoneNumber)
	params.SetTo(to)
	params.SetBody(message)

	_, err := client.Api.CreateMessage(params)
	return err
}

func (d *ItineraryDeliveryService) buildEmailBody(data *ItineraryData, fileURL, customMessage string) string {
	var body strings.Builder

	body.WriteString("<html><body>")
	body.WriteString("<h2>Your Travel Itinerary is Ready!</h2>")

	if customMessage != "" {
		body.WriteString(fmt.Sprintf("<p>%s</p>", customMessage))
	}

	body.WriteString(fmt.Sprintf("<p>Dear Traveler,</p>"))
	body.WriteString(fmt.Sprintf("<p>Your itinerary for <strong>%s</strong> is now ready for download.</p>", data.Destination))
	body.WriteString(fmt.Sprintf("<p><strong>Trip Details:</strong></p>"))
	body.WriteString("<ul>")
	body.WriteString(fmt.Sprintf("<li>Destination: %s</li>", data.Destination))
	body.WriteString(fmt.Sprintf("<li>Dates: %s to %s</li>",
		data.StartDate.Format("January 2, 2006"),
		data.EndDate.Format("January 2, 2006")))
	body.WriteString(fmt.Sprintf("<li>Travelers: %d</li>", data.Travelers))
	body.WriteString(fmt.Sprintf("<li>Total Budget: %.2f %s</li>", data.Budget, data.Currency))
	body.WriteString("</ul>")

	if fileURL != "" {
		body.WriteString(fmt.Sprintf("<p><a href=\"%s\" style=\"background-color: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;\">Download Your Itinerary</a></p>", fileURL))
	}

	body.WriteString("<p>Have a wonderful trip!</p>")
	body.WriteString("<p>Best regards,<br>AuraTravel AI Team</p>")
	body.WriteString("</body></html>")

	return body.String()
}

// Utility methods

func (d *ItineraryDeliveryService) getItineraryData(ctx context.Context, tripID, userID string) (*ItineraryData, error) {
	// Mock implementation - in production, fetch from database
	return &ItineraryData{
		TripID:      tripID,
		Destination: "Delhi, India",
		StartDate:   time.Now().AddDate(0, 0, 7),
		EndDate:     time.Now().AddDate(0, 0, 10),
		Travelers:   2,
		Budget:      50000,
		Currency:    "INR",
		Title:       "Amazing Delhi Adventure",
		Description: "A wonderful journey through India's capital",
		DailyItinerary: map[int]DayItinerary{
			1: {
				Date:      time.Now().AddDate(0, 0, 7),
				DayNumber: 1,
				Title:     "Arrival & Old Delhi Exploration",
				Morning: []Activity{
					{
						Name:        "Arrival at Delhi Airport",
						Type:        "transport",
						StartTime:   time.Now().AddDate(0, 0, 7).Add(8 * time.Hour),
						EndTime:     time.Now().AddDate(0, 0, 7).Add(9 * time.Hour),
						Location:    Location{Address: "Indira Gandhi International Airport"},
						Description: "Flight arrival and airport pickup",
						Cost:        0,
					},
				},
				Afternoon: []Activity{
					{
						Name:        "Red Fort",
						Type:        "sightseeing",
						StartTime:   time.Now().AddDate(0, 0, 7).Add(14 * time.Hour),
						EndTime:     time.Now().AddDate(0, 0, 7).Add(16 * time.Hour),
						Location:    Location{Address: "Red Fort, Old Delhi"},
						Description: "Historic Mughal fortress",
						Cost:        50,
					},
				},
				Evening: []Activity{
					{
						Name:        "Chandni Chowk Market",
						Type:        "shopping",
						StartTime:   time.Now().AddDate(0, 0, 7).Add(18 * time.Hour),
						EndTime:     time.Now().AddDate(0, 0, 7).Add(20 * time.Hour),
						Location:    Location{Address: "Chandni Chowk, Old Delhi"},
						Description: "Traditional market and street food",
						Cost:        200,
					},
				},
				TotalCost: 250,
			},
		},
		Hotels: []HotelBooking{
			{
				Name:            "Hotel New Delhi",
				Address:         "Connaught Place, New Delhi",
				CheckIn:         time.Now().AddDate(0, 0, 7),
				CheckOut:        time.Now().AddDate(0, 0, 10),
				RoomType:        "Deluxe Double",
				Nights:          3,
				TotalCost:       15000,
				ConfirmationNum: "HTL123456",
				Contact:         "+91-11-1234-5678",
				Amenities:       []string{"WiFi", "Breakfast", "Pool", "Gym"},
			},
		},
		Transportation: []TransportBooking{
			{
				Type:          "flight",
				From:          "Mumbai",
				To:            "Delhi",
				DepartureTime: time.Now().AddDate(0, 0, 7).Add(6 * time.Hour),
				ArrivalTime:   time.Now().AddDate(0, 0, 7).Add(8 * time.Hour),
				Provider:      "Air India",
				BookingRef:    "AI123456",
				Cost:          8000,
				Status:        "confirmed",
			},
		},
		ImportantInfo: []string{
			"Carry valid ID proof for all activities",
			"Weather might be hot - carry water and sun protection",
			"Respect local customs and dress codes",
			"Keep emergency contacts handy",
		},
		EmergencyContacts: []EmergencyContact{
			{
				Name:         "Emergency Services",
				Relationship: "emergency",
				Phone:        "112",
				Available24h: true,
			},
			{
				Name:         "Tourist Helpline",
				Relationship: "support",
				Phone:        "1363",
				Available24h: true,
			},
		},
		TotalCost:    25000,
		CreatedAt:    time.Now().AddDate(0, 0, -1),
		LastModified: time.Now(),
	}, nil
}

func (d *ItineraryDeliveryService) storeFile(ctx context.Context, fileData []byte, fileName, tripID string) (string, error) {
	// Mock file storage - in production, use cloud storage
	fileURL := fmt.Sprintf("%s/downloads/%s/%s", d.storageConfig.BaseURL, tripID, fileName)

	// In production, upload to cloud storage (Google Cloud Storage, AWS S3, etc.)
	log.Printf("Stored file %s for trip %s", fileName, tripID)

	return fileURL, nil
}

func (d *ItineraryDeliveryService) generateDeliveryID(tripID, userID string) string {
	return fmt.Sprintf("DEL_%s_%s_%d", tripID, userID, time.Now().Unix())
}

func (d *ItineraryDeliveryService) getUserEmail(ctx context.Context, userID string) (string, error) {
	// Mock implementation
	return "user@example.com", nil
}

func (d *ItineraryDeliveryService) getUserPhone(ctx context.Context, userID string) (string, error) {
	// Mock implementation
	return "+91-9876543210", nil
}

func (d *ItineraryDeliveryService) storeDeliveryRecord(ctx context.Context, result *DeliveryResult) {
	if d.firebase != nil {
		d.firebase.GetFirestoreClient().
			Collection("itinerary_deliveries").
			Doc(result.DeliveryID).
			Set(ctx, result)
	}
}

func (d *ItineraryDeliveryService) getDefaultHTMLTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}} - Travel Itinerary</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background-color: #007bff; color: white; padding: 20px; border-radius: 5px; }
        .day { margin: 20px 0; padding: 15px; border-left: 4px solid #007bff; }
        .activity { margin: 10px 0; padding: 10px; background-color: #f8f9fa; border-radius: 3px; }
        .time { font-weight: bold; color: #007bff; }
        .cost { color: #28a745; font-weight: bold; }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{.Title}}</h1>
        <p>{{.Destination}} | {{.StartDate.Format "January 2, 2006"}} - {{.EndDate.Format "January 2, 2006"}}</p>
    </div>
    
    <h2>Trip Overview</h2>
    <p><strong>Travelers:</strong> {{.Travelers}}</p>
    <p><strong>Budget:</strong> {{.Budget}} {{.Currency}}</p>
    <p><strong>Total Cost:</strong> <span class="cost">{{.TotalCost}} {{.Currency}}</span></p>
    
    <h2>Daily Itinerary</h2>
    {{range $dayNum, $day := .DailyItinerary}}
    <div class="day">
        <h3>Day {{$day.DayNumber}} - {{$day.Date.Format "Monday, January 2"}}</h3>
        {{range $day.Morning}}
        <div class="activity">
            <span class="time">{{.StartTime.Format "3:04 PM"}}</span> - {{.Name}}<br>
            <small>{{.Description}} | Cost: {{.Cost}} {{$.Currency}}</small>
        </div>
        {{end}}
        {{range $day.Afternoon}}
        <div class="activity">
            <span class="time">{{.StartTime.Format "3:04 PM"}}</span> - {{.Name}}<br>
            <small>{{.Description}} | Cost: {{.Cost}} {{$.Currency}}</small>
        </div>
        {{end}}
        {{range $day.Evening}}
        <div class="activity">
            <span class="time">{{.StartTime.Format "3:04 PM"}}</span> - {{.Name}}<br>
            <small>{{.Description}} | Cost: {{.Cost}} {{$.Currency}}</small>
        </div>
        {{end}}
    </div>
    {{end}}
    
    {{if .Hotels}}
    <h2>Accommodation</h2>
    {{range .Hotels}}
    <div class="activity">
        <strong>{{.Name}}</strong><br>
        {{.Address}}<br>
        Check-in: {{.CheckIn.Format "Jan 2, 2006"}} | Check-out: {{.CheckOut.Format "Jan 2, 2006"}}<br>
        {{if .ConfirmationNum}}Confirmation: {{.ConfirmationNum}}{{end}}
    </div>
    {{end}}
    {{end}}
    
    {{if .EmergencyContacts}}
    <h2>Emergency Contacts</h2>
    {{range .EmergencyContacts}}
    <p><strong>{{.Name}}</strong> ({{.Relationship}}): {{.Phone}}</p>
    {{end}}
    {{end}}
</body>
</html>
`
}

// GetDeliveryHistory retrieves delivery history for a trip
func (d *ItineraryDeliveryService) GetDeliveryHistory(ctx context.Context, tripID string) ([]*DeliveryResult, error) {
	if d.firebase == nil {
		return nil, fmt.Errorf("firebase service not available")
	}

	docs, err := d.firebase.GetFirestoreClient().
		Collection("itinerary_deliveries").
		Where("trip_id", "==", tripID).
		OrderBy("delivered_at", firestore.Desc).
		Documents(ctx).
		GetAll()

	if err != nil {
		return nil, err
	}

	var results []*DeliveryResult
	for _, doc := range docs {
		var result DeliveryResult
		if err := doc.DataTo(&result); err == nil {
			results = append(results, &result)
		}
	}

	return results, nil
}
