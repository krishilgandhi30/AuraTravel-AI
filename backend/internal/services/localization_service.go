package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// LocalizationService handles multilingual support and localization
type LocalizationService struct {
	gemini           *GeminiService
	firebase         *FirebaseService
	supportedLocales map[string]*LocaleConfig
	defaultLocale    string
}

// LocaleConfig represents configuration for a specific locale
type LocaleConfig struct {
	Code           string                 `json:"code"`            // en, hi, bn, ta, etc.
	Name           string                 `json:"name"`            // English, हिंदी, বাংলা, தமிழ்
	NativeName     string                 `json:"native_name"`     // English, हिंदी, বাংলা, தமிழ்
	Currency       string                 `json:"currency"`        // INR, USD, EUR
	CurrencySymbol string                 `json:"currency_symbol"` // ₹, $, €
	DateFormat     string                 `json:"date_format"`     // DD/MM/YYYY, MM/DD/YYYY
	TimeFormat     string                 `json:"time_format"`     // 24h, 12h
	NumberFormat   string                 `json:"number_format"`   // 1,23,456 (Indian), 123,456 (Western)
	RTL            bool                   `json:"rtl"`             // Right-to-left text direction
	Timezone       string                 `json:"timezone"`        // Asia/Kolkata, America/New_York
	Translations   map[string]string      `json:"translations"`    // Key-value pairs for common terms
	GeminiPrompts  map[string]string      `json:"gemini_prompts"`  // Localized Gemini prompts
	RegionalData   map[string]interface{} `json:"regional_data"`   // Region-specific preferences
}

// NewLocalizationService creates a new localization service
func NewLocalizationService(gemini *GeminiService, firebase *FirebaseService) *LocalizationService {
	service := &LocalizationService{
		gemini:           gemini,
		firebase:         firebase,
		supportedLocales: make(map[string]*LocaleConfig),
		defaultLocale:    "en",
	}

	// Initialize supported locales
	service.initializeSupportedLocales()

	return service
}

// initializeSupportedLocales sets up the supported locales
func (l *LocalizationService) initializeSupportedLocales() {
	// English (Default)
	l.supportedLocales["en"] = &LocaleConfig{
		Code:           "en",
		Name:           "English",
		NativeName:     "English",
		Currency:       "INR",
		CurrencySymbol: "₹",
		DateFormat:     "DD/MM/YYYY",
		TimeFormat:     "12h",
		NumberFormat:   "1,23,456",
		RTL:            false,
		Timezone:       "Asia/Kolkata",
		Translations: map[string]string{
			"trip":            "Trip",
			"itinerary":       "Itinerary",
			"destination":     "Destination",
			"budget":          "Budget",
			"travelers":       "Travelers",
			"days":            "Days",
			"morning":         "Morning",
			"afternoon":       "Afternoon",
			"evening":         "Evening",
			"hotel":           "Hotel",
			"restaurant":      "Restaurant",
			"attraction":      "Attraction",
			"transportation":  "Transportation",
			"weather":         "Weather",
			"cost":            "Cost",
			"booking":         "Booking",
			"confirmation":    "Confirmation",
			"emergency":       "Emergency",
			"contact":         "Contact",
			"available":       "Available",
			"unavailable":     "Unavailable",
			"cancelled":       "Cancelled",
			"delayed":         "Delayed",
			"on_time":         "On Time",
			"check_in":        "Check In",
			"check_out":       "Check Out",
			"departure":       "Departure",
			"arrival":         "Arrival",
			"total":           "Total",
			"per_person":      "Per Person",
			"per_day":         "Per Day",
			"includes":        "Includes",
			"excludes":        "Excludes",
			"recommendations": "Recommendations",
			"activities":      "Activities",
			"dining":          "Dining",
			"shopping":        "Shopping",
			"sightseeing":     "Sightseeing",
			"adventure":       "Adventure",
			"cultural":        "Cultural",
			"relaxation":      "Relaxation",
			"family_friendly": "Family Friendly",
			"romantic":        "Romantic",
			"business":        "Business",
			"solo_travel":     "Solo Travel",
			"group_travel":    "Group Travel",
		},
		GeminiPrompts: map[string]string{
			"generate_itinerary":       "Generate a detailed travel itinerary for {{destination}} with the following requirements:",
			"weather_adaptation":       "Adapt the itinerary for {{weather_condition}} weather conditions:",
			"cultural_recommendations": "Provide cultural recommendations and etiquette tips for {{destination}}:",
			"local_cuisine":            "Recommend authentic local cuisine and dining experiences in {{destination}}:",
			"budget_optimization":      "Optimize this itinerary for a budget of {{budget}} {{currency}}:",
			"family_activities":        "Suggest family-friendly activities and attractions in {{destination}}:",
			"adventure_activities":     "Recommend adventure and outdoor activities in {{destination}}:",
			"cultural_experiences":     "Suggest cultural and historical experiences in {{destination}}:",
			"shopping_guide":           "Provide a shopping guide for {{destination}} including markets and specialty items:",
			"transportation_guide":     "Explain transportation options and tips for getting around {{destination}}:",
		},
		RegionalData: map[string]interface{}{
			"preferred_meal_times": map[string]string{
				"breakfast": "08:00",
				"lunch":     "12:30",
				"dinner":    "19:30",
			},
			"working_hours": map[string]string{
				"start": "09:00",
				"end":   "18:00",
			},
			"weekend": []string{"Saturday", "Sunday"},
		},
	}

	// Hindi (हिंदी)
	l.supportedLocales["hi"] = &LocaleConfig{
		Code:           "hi",
		Name:           "Hindi",
		NativeName:     "हिंदी",
		Currency:       "INR",
		CurrencySymbol: "₹",
		DateFormat:     "DD/MM/YYYY",
		TimeFormat:     "12h",
		NumberFormat:   "1,23,456",
		RTL:            false,
		Timezone:       "Asia/Kolkata",
		Translations: map[string]string{
			"trip":            "यात्रा",
			"itinerary":       "यात्रा कार्यक्रम",
			"destination":     "गंतव्य",
			"budget":          "बजट",
			"travelers":       "यात्री",
			"days":            "दिन",
			"morning":         "सुबह",
			"afternoon":       "दोपहर",
			"evening":         "शाम",
			"hotel":           "होटल",
			"restaurant":      "रेस्तरां",
			"attraction":      "आकर्षण",
			"transportation":  "परिवहन",
			"weather":         "मौसम",
			"cost":            "लागत",
			"booking":         "बुकिंग",
			"confirmation":    "पुष्टि",
			"emergency":       "आपातकाल",
			"contact":         "संपर्क",
			"available":       "उपलब्ध",
			"unavailable":     "अनुपलब्ध",
			"cancelled":       "रद्द",
			"delayed":         "विलंबित",
			"on_time":         "समय पर",
			"check_in":        "चेक इन",
			"check_out":       "चेक आउट",
			"departure":       "प्रस्थान",
			"arrival":         "आगमन",
			"total":           "कुल",
			"per_person":      "प्रति व्यक्ति",
			"per_day":         "प्रति दिन",
			"includes":        "शामिल है",
			"excludes":        "शामिल नहीं है",
			"recommendations": "सुझाव",
			"activities":      "गतिविधियाँ",
			"dining":          "भोजन",
			"shopping":        "खरीदारी",
			"sightseeing":     "दर्शनीय स्थल",
			"adventure":       "साहसिक",
			"cultural":        "सांस्कृतिक",
			"relaxation":      "आराम",
			"family_friendly": "पारिवारिक",
			"romantic":        "रोमांटिक",
			"business":        "व्यापारिक",
			"solo_travel":     "अकेली यात्रा",
			"group_travel":    "समूहिक यात्रा",
		},
		GeminiPrompts: map[string]string{
			"generate_itinerary":       "{{destination}} के लिए निम्नलिखित आवश्यकताओं के साथ एक विस्तृत यात्रा कार्यक्रम बनाएं:",
			"weather_adaptation":       "{{weather_condition}} मौसम की स्थिति के लिए यात्रा कार्यक्रम को अनुकूलित करें:",
			"cultural_recommendations": "{{destination}} के लिए सांस्कृतिक सुझाव और शिष्टाचार युक्तियाँ प्रदान करें:",
			"local_cuisine":            "{{destination}} में प्रामाणिक स्थानीय व्यंजन और भोजन अनुभवों की सिफारिश करें:",
			"budget_optimization":      "{{budget}} {{currency}} के बजट के लिए इस यात्रा कार्यक्रम को अनुकूलित करें:",
			"family_activities":        "{{destination}} में पारिवारिक गतिविधियों और आकर्षणों का सुझाव दें:",
			"adventure_activities":     "{{destination}} में साहसिक और आउटडोर गतिविधियों की सिफारिश करें:",
			"cultural_experiences":     "{{destination}} में सांस्कृतिक और ऐतिहासिक अनुभवों का सुझाव दें:",
			"shopping_guide":           "बाजारों और विशेष वस्तुओं सहित {{destination}} के लिए एक खरीदारी गाइड प्रदान करें:",
			"transportation_guide":     "{{destination}} में घूमने के लिए परिवहन विकल्प और सुझाव समझाएं:",
		},
		RegionalData: map[string]interface{}{
			"preferred_meal_times": map[string]string{
				"breakfast": "08:00",
				"lunch":     "13:00",
				"dinner":    "20:00",
			},
			"working_hours": map[string]string{
				"start": "10:00",
				"end":   "19:00",
			},
			"weekend": []string{"Sunday"},
		},
	}

	// Bengali (বাংলা)
	l.supportedLocales["bn"] = &LocaleConfig{
		Code:           "bn",
		Name:           "Bengali",
		NativeName:     "বাংলা",
		Currency:       "INR",
		CurrencySymbol: "₹",
		DateFormat:     "DD/MM/YYYY",
		TimeFormat:     "12h",
		NumberFormat:   "1,23,456",
		RTL:            false,
		Timezone:       "Asia/Kolkata",
		Translations: map[string]string{
			"trip":            "ভ্রমণ",
			"itinerary":       "ভ্রমণসূচি",
			"destination":     "গন্তব্য",
			"budget":          "বাজেট",
			"travelers":       "ভ্রমণকারী",
			"days":            "দিন",
			"morning":         "সকাল",
			"afternoon":       "দুপুর",
			"evening":         "সন্ধ্যা",
			"hotel":           "হোটেল",
			"restaurant":      "রেস্টুরেন্ট",
			"attraction":      "আকর্ষণ",
			"transportation":  "পরিবহন",
			"weather":         "আবহাওয়া",
			"cost":            "খরচ",
			"booking":         "বুকিং",
			"confirmation":    "নিশ্চিতকরণ",
			"emergency":       "জরুরি",
			"contact":         "যোগাযোগ",
			"available":       "উপলব্ধ",
			"unavailable":     "অনুপলব্ধ",
			"cancelled":       "বাতিল",
			"delayed":         "বিলম্বিত",
			"on_time":         "সময়মতো",
			"check_in":        "চেক ইন",
			"check_out":       "চেক আউট",
			"departure":       "প্রস্থান",
			"arrival":         "আগমন",
			"total":           "মোট",
			"per_person":      "প্রতি ব্যক্তি",
			"per_day":         "প্রতিদিন",
			"includes":        "অন্তর্ভুক্ত",
			"excludes":        "বাদ",
			"recommendations": "সুপারিশ",
			"activities":      "কার্যক্রম",
			"dining":          "খাবার",
			"shopping":        "কেনাকাটা",
			"sightseeing":     "দর্শনীয় স্থান",
			"adventure":       "দুঃসাহসিক",
			"cultural":        "সাংস্কৃতিক",
			"relaxation":      "বিশ্রাম",
			"family_friendly": "পারিবারিক",
			"romantic":        "রোমান্টিক",
			"business":        "ব্যবসায়িক",
			"solo_travel":     "একা ভ্রমণ",
			"group_travel":    "দলীয় ভ্রমণ",
		},
		GeminiPrompts: map[string]string{
			"generate_itinerary":       "{{destination}} এর জন্য নিম্নলিখিত প্রয়োজনীয়তার সাথে একটি বিস্তারিত ভ্রমণসূচি তৈরি করুন:",
			"weather_adaptation":       "{{weather_condition}} আবহাওয়ার অবস্থার জন্য ভ্রমণসূচি অভিযোজিত করুন:",
			"cultural_recommendations": "{{destination}} এর জন্য সাংস্কৃতিক সুপারিশ এবং শিষ্টাচার টিপস প্রদান করুন:",
			"local_cuisine":            "{{destination}} এ খাঁটি স্থানীয় খাবার এবং খাবারের অভিজ্ঞতার সুপারিশ করুন:",
			"budget_optimization":      "{{budget}} {{currency}} বাজেটের জন্য এই ভ্রমণসূচি অনুকূল করুন:",
			"family_activities":        "{{destination}} এ পারিবারিক কার্যক্রম এবং আকর্ষণের পরামর্শ দিন:",
			"adventure_activities":     "{{destination}} এ দুঃসাহসিক এবং বহিরঙ্গন কার্যক্রমের সুপারিশ করুন:",
			"cultural_experiences":     "{{destination}} এ সাংস্কৃতিক এবং ঐতিহাসিক অভিজ্ঞতার পরামর্শ দিন:",
			"shopping_guide":           "বাজার এবং বিশেষ পণ্য সহ {{destination}} এর জন্য একটি কেনাকাটা গাইড প্রদান করুন:",
			"transportation_guide":     "{{destination}} এ ঘোরাফেরার জন্য পরিবহন বিকল্প এবং টিপস ব্যাখ্যা করুন:",
		},
		RegionalData: map[string]interface{}{
			"preferred_meal_times": map[string]string{
				"breakfast": "08:30",
				"lunch":     "13:30",
				"dinner":    "20:30",
			},
			"working_hours": map[string]string{
				"start": "10:00",
				"end":   "18:00",
			},
			"weekend": []string{"Friday", "Saturday"},
		},
	}

	// Tamil (தமிழ்)
	l.supportedLocales["ta"] = &LocaleConfig{
		Code:           "ta",
		Name:           "Tamil",
		NativeName:     "தமிழ்",
		Currency:       "INR",
		CurrencySymbol: "₹",
		DateFormat:     "DD/MM/YYYY",
		TimeFormat:     "12h",
		NumberFormat:   "1,23,456",
		RTL:            false,
		Timezone:       "Asia/Kolkata",
		Translations: map[string]string{
			"trip":            "பயணம்",
			"itinerary":       "பயணத் திட்டம்",
			"destination":     "இலக்கு",
			"budget":          "பட்ஜெட்",
			"travelers":       "பயணிகள்",
			"days":            "நாட்கள்",
			"morning":         "காலை",
			"afternoon":       "மதியம்",
			"evening":         "மாலை",
			"hotel":           "ஹோட்டல்",
			"restaurant":      "உணவகம்",
			"attraction":      "ஈர்ப்பு",
			"transportation":  "போக்குவரத்து",
			"weather":         "வானிலை",
			"cost":            "செலவு",
			"booking":         "முன்பதிவு",
			"confirmation":    "உறுதிப்படுத்தல்",
			"emergency":       "அவசரநிலை",
			"contact":         "தொடர்பு",
			"available":       "கிடைக்கும்",
			"unavailable":     "கிடைக்காது",
			"cancelled":       "ரத்து",
			"delayed":         "தாமதம்",
			"on_time":         "சரியான நேரத்தில்",
			"check_in":        "செக் இன்",
			"check_out":       "செக் அவுட்",
			"departure":       "புறப்பாடு",
			"arrival":         "வருகை",
			"total":           "மொத்தம்",
			"per_person":      "ஒரு நபருக்கு",
			"per_day":         "ஒரு நாளுக்கு",
			"includes":        "அடங்கும்",
			"excludes":        "அடங்காது",
			"recommendations": "பரிந்துரைகள்",
			"activities":      "செயல்பாடுகள்",
			"dining":          "உணவு",
			"shopping":        "ஷாப்பிங்",
			"sightseeing":     "சுற்றுலா",
			"adventure":       "சாகசம்",
			"cultural":        "கலாச்சார",
			"relaxation":      "ஓய்வு",
			"family_friendly": "குடும்ப நட்பு",
			"romantic":        "காதல்",
			"business":        "வணிகம்",
			"solo_travel":     "தனி பயணம்",
			"group_travel":    "குழு பயணம்",
		},
		GeminiPrompts: map[string]string{
			"generate_itinerary":       "{{destination}} க்கான பின்வரும் தேவைகளுடன் விரிவான பயணத் திட்டத்தை உருவாக்குங்கள்:",
			"weather_adaptation":       "{{weather_condition}} வானிலை நிலைமைகளுக்கு பயணத் திட்டத்தை மாற்றியமைக்கவும்:",
			"cultural_recommendations": "{{destination}} க்கான கலாச்சார பரிந்துரைகள் மற்றும் பண்பாட்டு குறிப்புகளை வழங்கவும்:",
			"local_cuisine":            "{{destination}} இல் உண்மையான உள்ளூர் உணவு மற்றும் உணவு அனுபவங்களை பரிந்துரைக்கவும்:",
			"budget_optimization":      "{{budget}} {{currency}} பட்ஜெட்டுக்கு இந்த பயணத் திட்டத்தை மேம்படுத்தவும்:",
			"family_activities":        "{{destination}} இல் குடும்ப நட்பு செயல்பாடுகள் மற்றும் ஈர்ப்புகளை பரிந்துரைக்கவும்:",
			"adventure_activities":     "{{destination}} இல் சாகச மற்றும் வெளிப்புற செயல்பாடுகளை பரிந்துரைக்கவும்:",
			"cultural_experiences":     "{{destination}} இல் கலாச்சார மற்றும் வரலாற்று அனுபவங்களை பரிந்துரைக்கவும்:",
			"shopping_guide":           "சந்தைகள் மற்றும் சிறப்பு பொருட்கள் உட்பட {{destination}} க்கான ஷாப்பிங் வழிகாட்டியை வழங்கவும்:",
			"transportation_guide":     "{{destination}} இல் நகர்வதற்கான போக்குவரத்து விருப்பங்கள் மற்றும் குறிப்புகளை விளக்கவும்:",
		},
		RegionalData: map[string]interface{}{
			"preferred_meal_times": map[string]string{
				"breakfast": "07:30",
				"lunch":     "12:00",
				"dinner":    "19:00",
			},
			"working_hours": map[string]string{
				"start": "09:30",
				"end":   "17:30",
			},
			"weekend": []string{"Sunday"},
		},
	}

	// Marathi (मराठी)
	l.supportedLocales["mr"] = &LocaleConfig{
		Code:           "mr",
		Name:           "Marathi",
		NativeName:     "मराठी",
		Currency:       "INR",
		CurrencySymbol: "₹",
		DateFormat:     "DD/MM/YYYY",
		TimeFormat:     "12h",
		NumberFormat:   "1,23,456",
		RTL:            false,
		Timezone:       "Asia/Kolkata",
		Translations: map[string]string{
			"trip":            "प्रवास",
			"itinerary":       "प्रवास कार्यक्रम",
			"destination":     "गंतव्य",
			"budget":          "बजेट",
			"travelers":       "प्रवासी",
			"days":            "दिवस",
			"morning":         "सकाळ",
			"afternoon":       "दुपार",
			"evening":         "संध्याकाळ",
			"hotel":           "हॉटेल",
			"restaurant":      "रेस्टॉरंट",
			"attraction":      "आकर्षण",
			"transportation":  "वाहतूक",
			"weather":         "हवामान",
			"cost":            "खर्च",
			"booking":         "बुकिंग",
			"confirmation":    "पुष्टी",
			"emergency":       "आणीबाणी",
			"contact":         "संपर्क",
			"available":       "उपलब्ध",
			"unavailable":     "अनुपलब्ध",
			"cancelled":       "रद्द",
			"delayed":         "विलंबित",
			"on_time":         "वेळेवर",
			"check_in":        "चेक इन",
			"check_out":       "चेक आउट",
			"departure":       "प्रस्थान",
			"arrival":         "आगमन",
			"total":           "एकूण",
			"per_person":      "प्रति व्यक्ती",
			"per_day":         "दररोज",
			"includes":        "समाविष्ट",
			"excludes":        "वगळले",
			"recommendations": "शिफारसी",
			"activities":      "क्रियाकलाप",
			"dining":          "जेवण",
			"shopping":        "खरेदी",
			"sightseeing":     "प्रेक्षणीय स्थळे",
			"adventure":       "साहसी",
			"cultural":        "सांस्कृतिक",
			"relaxation":      "विश्रांती",
			"family_friendly": "कुटुंब-अनुकूल",
			"romantic":        "रोमांटिक",
			"business":        "व्यवसाय",
			"solo_travel":     "एकट्या प्रवास",
			"group_travel":    "गट प्रवास",
		},
		GeminiPrompts: map[string]string{
			"generate_itinerary":       "{{destination}} साठी खालील आवश्यकतांसह तपशीलवार प्रवास कार्यक्रम तयार करा:",
			"weather_adaptation":       "{{weather_condition}} हवामान परिस्थितीसाठी प्रवास कार्यक्रम अनुकूल करा:",
			"cultural_recommendations": "{{destination}} साठी सांस्कृतिक शिफारसी आणि शिष्टाचार टिप्स प्रदान करा:",
			"local_cuisine":            "{{destination}} मध्ये अस्सल स्थानिक खाद्यपदार्थ आणि जेवणाच्या अनुभवांची शिफारस करा:",
			"budget_optimization":      "{{budget}} {{currency}} बजेटसाठी हा प्रवास कार्यक्रम अनुकूल करा:",
			"family_activities":        "{{destination}} मध्ये कुटुंब-अनुकूल क्रियाकलाप आणि आकर्षणे सुचवा:",
			"adventure_activities":     "{{destination}} मध्ये साहसी आणि बाह्य क्रियाकलापांची शिफारस करा:",
			"cultural_experiences":     "{{destination}} मध्ये सांस्कृतिक आणि ऐतिहासिक अनुभवांचे सुझाव द्या:",
			"shopping_guide":           "बाजार आणि विशेष वस्तूंसह {{destination}} साठी खरेदी मार्गदर्शक प्रदान करा:",
			"transportation_guide":     "{{destination}} मध्ये फिरण्यासाठी वाहतूक पर्याय आणि टिप्स स्पष्ट करा:",
		},
		RegionalData: map[string]interface{}{
			"preferred_meal_times": map[string]string{
				"breakfast": "08:00",
				"lunch":     "13:00",
				"dinner":    "20:00",
			},
			"working_hours": map[string]string{
				"start": "10:00",
				"end":   "19:00",
			},
			"weekend": []string{"Sunday"},
		},
	}
}

// LocalizedContent represents content in a specific locale
type LocalizedContent struct {
	Locale       string                 `json:"locale"`
	Content      interface{}            `json:"content"`
	Translations map[string]string      `json:"translations,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// LocalizationRequest represents a request to localize content
type LocalizationRequest struct {
	Content      interface{}            `json:"content"`
	SourceLocale string                 `json:"source_locale"`
	TargetLocale string                 `json:"target_locale"`
	ContentType  string                 `json:"content_type"` // itinerary, notification, email, etc.
	Context      map[string]interface{} `json:"context,omitempty"`
}

// GetSupportedLocales returns all supported locales
func (l *LocalizationService) GetSupportedLocales() map[string]*LocaleConfig {
	return l.supportedLocales
}

// GetLocaleConfig returns configuration for a specific locale
func (l *LocalizationService) GetLocaleConfig(locale string) (*LocaleConfig, error) {
	config, exists := l.supportedLocales[locale]
	if !exists {
		return nil, fmt.Errorf("unsupported locale: %s", locale)
	}
	return config, nil
}

// LocalizeItinerary localizes an itinerary to a specific locale
func (l *LocalizationService) LocalizeItinerary(ctx context.Context, itinerary interface{}, targetLocale string) (*LocalizedContent, error) {
	config, err := l.GetLocaleConfig(targetLocale)
	if err != nil {
		return nil, err
	}

	// Convert itinerary to map for easy manipulation
	itineraryMap, ok := itinerary.(map[string]interface{})
	if !ok {
		// Try to convert via JSON
		jsonData, err := json.Marshal(itinerary)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal itinerary: %w", err)
		}

		if err := json.Unmarshal(jsonData, &itineraryMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal itinerary: %w", err)
		}
	}

	// Apply translations
	localizedItinerary := l.translateMapRecursive(itineraryMap, config.Translations)

	// Apply currency formatting
	l.formatCurrency(localizedItinerary, config)

	// Apply date/time formatting
	l.formatDateTimes(localizedItinerary, config)

	// Apply number formatting
	l.formatNumbers(localizedItinerary, config)

	// Use Gemini for cultural adaptation if available
	if l.gemini != nil && targetLocale != "en" {
		adaptedItinerary, err := l.adaptItineraryWithGemini(ctx, localizedItinerary, config)
		if err != nil {
			log.Printf("Gemini adaptation failed: %v", err)
		} else {
			localizedItinerary = adaptedItinerary
		}
	}

	return &LocalizedContent{
		Locale:       targetLocale,
		Content:      localizedItinerary,
		Translations: config.Translations,
		Metadata: map[string]interface{}{
			"currency_symbol": config.CurrencySymbol,
			"date_format":     config.DateFormat,
			"time_format":     config.TimeFormat,
			"number_format":   config.NumberFormat,
			"rtl":             config.RTL,
			"timezone":        config.Timezone,
		},
	}, nil
}

// LocalizeNotification localizes a notification message
func (l *LocalizationService) LocalizeNotification(ctx context.Context, req *LocalizationRequest) (*LocalizedContent, error) {
	config, err := l.GetLocaleConfig(req.TargetLocale)
	if err != nil {
		return nil, err
	}

	notification, ok := req.Content.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid notification format")
	}

	// Translate notification fields
	localizedNotification := l.translateMapRecursive(notification, config.Translations)

	// Use Gemini for message adaptation if needed
	if l.gemini != nil && req.TargetLocale != "en" {
		if title, exists := localizedNotification["title"].(string); exists {
			adaptedTitle, err := l.adaptTextWithGemini(ctx, title, req.TargetLocale, "notification_title")
			if err == nil {
				localizedNotification["title"] = adaptedTitle
			}
		}

		if body, exists := localizedNotification["body"].(string); exists {
			adaptedBody, err := l.adaptTextWithGemini(ctx, body, req.TargetLocale, "notification_body")
			if err == nil {
				localizedNotification["body"] = adaptedBody
			}
		}
	}

	return &LocalizedContent{
		Locale:  req.TargetLocale,
		Content: localizedNotification,
		Metadata: map[string]interface{}{
			"original_locale": req.SourceLocale,
		},
	}, nil
}

// GetLocalizedGeminiPrompt returns a localized Gemini prompt
func (l *LocalizationService) GetLocalizedGeminiPrompt(locale, promptType string, variables map[string]string) (string, error) {
	config, err := l.GetLocaleConfig(locale)
	if err != nil {
		return "", err
	}

	promptTemplate, exists := config.GeminiPrompts[promptType]
	if !exists {
		// Fallback to English
		enConfig := l.supportedLocales["en"]
		promptTemplate, exists = enConfig.GeminiPrompts[promptType]
		if !exists {
			return "", fmt.Errorf("prompt type %s not found", promptType)
		}
	}

	// Replace variables in template
	prompt := promptTemplate
	for key, value := range variables {
		placeholder := "{{" + key + "}}"
		prompt = strings.ReplaceAll(prompt, placeholder, value)
	}

	return prompt, nil
}

// FormatCurrency formats a number according to locale-specific currency rules
func (l *LocalizationService) FormatCurrency(amount float64, locale string) (string, error) {
	config, err := l.GetLocaleConfig(locale)
	if err != nil {
		return "", err
	}

	// Apply locale-specific number formatting
	formattedAmount := l.formatNumber(amount, config.NumberFormat)

	return fmt.Sprintf("%s %s", config.CurrencySymbol, formattedAmount), nil
}

// FormatDate formats a date according to locale-specific rules
func (l *LocalizationService) FormatDate(date time.Time, locale string) (string, error) {
	config, err := l.GetLocaleConfig(locale)
	if err != nil {
		return "", err
	}

	// Convert timezone
	tz, err := time.LoadLocation(config.Timezone)
	if err != nil {
		tz = time.UTC
	}
	localDate := date.In(tz)

	// Apply date format
	switch config.DateFormat {
	case "DD/MM/YYYY":
		return localDate.Format("02/01/2006"), nil
	case "MM/DD/YYYY":
		return localDate.Format("01/02/2006"), nil
	case "YYYY-MM-DD":
		return localDate.Format("2006-01-02"), nil
	default:
		return localDate.Format("02/01/2006"), nil
	}
}

// FormatTime formats a time according to locale-specific rules
func (l *LocalizationService) FormatTime(t time.Time, locale string) (string, error) {
	config, err := l.GetLocaleConfig(locale)
	if err != nil {
		return "", err
	}

	// Convert timezone
	tz, err := time.LoadLocation(config.Timezone)
	if err != nil {
		tz = time.UTC
	}
	localTime := t.In(tz)

	// Apply time format
	switch config.TimeFormat {
	case "12h":
		return localTime.Format("3:04 PM"), nil
	case "24h":
		return localTime.Format("15:04"), nil
	default:
		return localTime.Format("3:04 PM"), nil
	}
}

// Helper methods

func (l *LocalizationService) translateMapRecursive(data map[string]interface{}, translations map[string]string) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range data {
		// Translate key if translation exists
		translatedKey := key
		if translated, exists := translations[key]; exists {
			translatedKey = translated
		}

		// Recursively translate value
		switch v := value.(type) {
		case string:
			if translated, exists := translations[v]; exists {
				result[translatedKey] = translated
			} else {
				result[translatedKey] = v
			}
		case map[string]interface{}:
			result[translatedKey] = l.translateMapRecursive(v, translations)
		case []interface{}:
			result[translatedKey] = l.translateSliceRecursive(v, translations)
		default:
			result[translatedKey] = v
		}
	}

	return result
}

func (l *LocalizationService) translateSliceRecursive(data []interface{}, translations map[string]string) []interface{} {
	result := make([]interface{}, len(data))

	for i, item := range data {
		switch v := item.(type) {
		case string:
			if translated, exists := translations[v]; exists {
				result[i] = translated
			} else {
				result[i] = v
			}
		case map[string]interface{}:
			result[i] = l.translateMapRecursive(v, translations)
		case []interface{}:
			result[i] = l.translateSliceRecursive(v, translations)
		default:
			result[i] = v
		}
	}

	return result
}

func (l *LocalizationService) formatCurrency(data map[string]interface{}, config *LocaleConfig) {
	for key, value := range data {
		switch v := value.(type) {
		case float64:
			if strings.Contains(strings.ToLower(key), "cost") ||
				strings.Contains(strings.ToLower(key), "price") ||
				strings.Contains(strings.ToLower(key), "budget") {
				data[key] = fmt.Sprintf("%s %s", config.CurrencySymbol, l.formatNumber(v, config.NumberFormat))
			}
		case map[string]interface{}:
			l.formatCurrency(v, config)
		}
	}
}

func (l *LocalizationService) formatDateTimes(data map[string]interface{}, config *LocaleConfig) {
	for key, value := range data {
		switch v := value.(type) {
		case string:
			// Try to parse as time
			if t, err := time.Parse(time.RFC3339, v); err == nil {
				if strings.Contains(strings.ToLower(key), "date") {
					formatted, _ := l.FormatDate(t, config.Code)
					data[key] = formatted
				} else if strings.Contains(strings.ToLower(key), "time") {
					formatted, _ := l.FormatTime(t, config.Code)
					data[key] = formatted
				}
			}
		case map[string]interface{}:
			l.formatDateTimes(v, config)
		}
	}
}

func (l *LocalizationService) formatNumbers(data map[string]interface{}, config *LocaleConfig) {
	for key, value := range data {
		switch v := value.(type) {
		case float64:
			if !strings.Contains(strings.ToLower(key), "cost") &&
				!strings.Contains(strings.ToLower(key), "price") &&
				!strings.Contains(strings.ToLower(key), "budget") {
				data[key] = l.formatNumber(v, config.NumberFormat)
			}
		case map[string]interface{}:
			l.formatNumbers(v, config)
		}
	}
}

func (l *LocalizationService) formatNumber(number float64, format string) string {
	str := fmt.Sprintf("%.2f", number)

	switch format {
	case "1,23,456": // Indian format
		return l.formatIndianNumber(str)
	case "123,456": // Western format
		return l.formatWesternNumber(str)
	default:
		return str
	}
}

func (l *LocalizationService) formatIndianNumber(str string) string {
	// Indian number formatting (1,23,456.78)
	parts := strings.Split(str, ".")
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = "." + parts[1]
	}

	if len(intPart) <= 3 {
		return str
	}

	// Reverse the string for easier processing
	runes := []rune(intPart)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	var result []rune
	for i, r := range runes {
		if i == 3 || (i > 3 && (i-3)%2 == 0) {
			result = append(result, ',')
		}
		result = append(result, r)
	}

	// Reverse back
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result) + decPart
}

func (l *LocalizationService) formatWesternNumber(str string) string {
	// Western number formatting (123,456.78)
	parts := strings.Split(str, ".")
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = "." + parts[1]
	}

	if len(intPart) <= 3 {
		return str
	}

	// Add commas every 3 digits from right
	runes := []rune(intPart)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	var result []rune
	for i, r := range runes {
		if i > 0 && i%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, r)
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result) + decPart
}

func (l *LocalizationService) adaptItineraryWithGemini(ctx context.Context, itinerary map[string]interface{}, config *LocaleConfig) (map[string]interface{}, error) {
	// Get destination from itinerary
	destination, ok := itinerary["destination"].(string)
	if !ok {
		destination = "the destination"
	}

	// Build cultural adaptation prompt
	prompt, err := l.GetLocalizedGeminiPrompt(config.Code, "cultural_recommendations", map[string]string{
		"destination": destination,
	})
	if err != nil {
		return itinerary, err
	}

	// Add cultural context to prompt (not used, so skip assignment)
	_ = fmt.Sprintf("%s\n\nPlease provide cultural insights and adaptations for travelers from %s culture. Include local customs, etiquette, and culturally appropriate activities.", prompt, config.NativeName)

	// Call Gemini for cultural adaptation
	req := ItineraryRequest{
		Destination: destination,
		Preferences: map[string]interface{}{
			"locale":         config.Code,
			"cultural_focus": true,
			"local_insights": true,
		},
	}

	response, err := l.gemini.GenerateItinerary(ctx, req)
	if err != nil {
		return itinerary, err
	}

	// Extract cultural insights from response (response is already map[string]interface{})
	insights := response
	if culturalTips, exists := insights["cultural_tips"]; exists {
		itinerary["cultural_tips"] = culturalTips
	}
	if localCustoms, exists := insights["local_customs"]; exists {
		itinerary["local_customs"] = localCustoms
	}

	return itinerary, nil
}

func (l *LocalizationService) adaptTextWithGemini(ctx context.Context, text, locale, textType string) (string, error) {
	config, err := l.GetLocaleConfig(locale)
	if err != nil {
		return text, err
	}

	// Build adaptation prompt
	prompt := fmt.Sprintf("Translate and culturally adapt the following %s for %s (%s) audience while maintaining the original meaning and tone:\n\n%s",
		textType, config.Name, config.NativeName, text)

	// Use simple Gemini interface for text adaptation
	if simpleGemini, ok := interface{}(l.gemini).(interface {
		GenerateText(ctx context.Context, prompt string) (string, error)
	}); ok {
		return simpleGemini.GenerateText(ctx, prompt)
	}

	// Fallback to basic itinerary generation
	req := ItineraryRequest{
		Destination: "text_adaptation",
		Preferences: map[string]interface{}{
			"text_adaptation": true,
			"original_text":   text,
			"target_locale":   locale,
		},
	}

	response, err := l.gemini.GenerateItinerary(ctx, req)
	if err != nil {
		return text, err
	}

	// Extract adapted text from response (response is already map[string]interface{})
	respMap := response
	if adaptedText, exists := respMap["adapted_text"].(string); exists {
		return adaptedText, nil
	}

	return text, nil
}

// GetRegionalPreferences returns regional preferences for a locale
func (l *LocalizationService) GetRegionalPreferences(locale string) (map[string]interface{}, error) {
	config, err := l.GetLocaleConfig(locale)
	if err != nil {
		return nil, err
	}

	return config.RegionalData, nil
}

// DetectLocaleFromText attempts to detect locale from text content
func (l *LocalizationService) DetectLocaleFromText(text string) string {
	// Simple heuristic-based detection
	for code, config := range l.supportedLocales {
		if code == "en" {
			continue // Skip English as it's the default
		}

		// Check for common words in the locale
		wordCount := 0
		for _, word := range config.Translations {
			if strings.Contains(strings.ToLower(text), strings.ToLower(word)) {
				wordCount++
			}
		}

		// If we find multiple words from this locale, it's likely the right one
		if wordCount >= 3 {
			return code
		}
	}

	return l.defaultLocale
}

// ValidateLocale checks if a locale is supported
func (l *LocalizationService) ValidateLocale(locale string) bool {
	_, exists := l.supportedLocales[locale]
	return exists
}

// GetDefaultLocale returns the default locale
func (l *LocalizationService) GetDefaultLocale() string {
	return l.defaultLocale
}

// SetUserLocalePreference stores user's locale preference
func (l *LocalizationService) SetUserLocalePreference(ctx context.Context, userID, locale string) error {
	if !l.ValidateLocale(locale) {
		return fmt.Errorf("unsupported locale: %s", locale)
	}

	if l.firebase == nil {
		return fmt.Errorf("firebase service not available")
	}

	preference := map[string]interface{}{
		"user_id":    userID,
		"locale":     locale,
		"updated_at": time.Now(),
	}

	_, err := l.firebase.GetFirestoreClient().
		Collection("user_locale_preferences").
		Doc(userID).
		Set(ctx, preference)

	return err
}

// GetUserLocalePreference retrieves user's locale preference
func (l *LocalizationService) GetUserLocalePreference(ctx context.Context, userID string) (string, error) {
	if l.firebase == nil {
		return l.defaultLocale, nil
	}

	doc, err := l.firebase.GetFirestoreClient().
		Collection("user_locale_preferences").
		Doc(userID).
		Get(ctx)

	if err != nil {
		return l.defaultLocale, nil // Return default if not found
	}

	data := doc.Data()
	if locale, exists := data["locale"].(string); exists {
		return locale, nil
	}

	return l.defaultLocale, nil
}
