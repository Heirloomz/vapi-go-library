package chat

// ChatMessage represents a message in a chat conversation
type ChatMessage struct {
	Role             string `json:"role"`
	Content          string `json:"content"`
	Time             int64  `json:"time"`
	SecondsFromStart int64  `json:"secondsFromStart"`
}

// ChatInput represents the input for creating a chat
// Can be either a string or an array of ChatMessage objects
type ChatInput interface{}

// CreateChatRequest represents the request payload for creating a chat
type CreateChatRequest struct {
	// Required: Input text or array of chat messages
	Input ChatInput `json:"input"`

	// Optional: Assistant configuration
	AssistantID        *string             `json:"assistantId,omitempty"`
	Assistant          *Assistant          `json:"assistant,omitempty"`
	AssistantOverrides *AssistantOverrides `json:"assistantOverrides,omitempty"`

	// Optional: Chat metadata
	Name           *string `json:"name,omitempty"`
	SessionID      *string `json:"sessionId,omitempty"`
	PreviousChatID *string `json:"previousChatId,omitempty"`

	// Optional: Streaming configuration
	Stream *bool `json:"stream,omitempty"`
}

// ChatResponse represents the response from creating a chat
type ChatResponse struct {
	ID             string        `json:"id"`
	OrgID          string        `json:"orgId"`
	AssistantID    *string       `json:"assistantId,omitempty"`
	Assistant      *Assistant    `json:"assistant,omitempty"`
	Name           *string       `json:"name,omitempty"`
	SessionID      *string       `json:"sessionId,omitempty"`
	Input          ChatInput     `json:"input"`
	Messages       []ChatMessage `json:"messages"`
	Output         []ChatMessage `json:"output"`
	Stream         *bool         `json:"stream,omitempty"`
	PreviousChatID *string       `json:"previousChatId,omitempty"`
	CreatedAt      string        `json:"createdAt"`
	UpdatedAt      string        `json:"updatedAt"`
	Costs          []Cost        `json:"costs"`
	Cost           float64       `json:"cost"`
}

// StreamingChatResponse represents a streaming chat response
type StreamingChatResponse struct {
	ID      string `json:"id"`
	OrgID   string `json:"orgId"`
	Message string `json:"message"`
	Done    bool   `json:"done"`
}

// Cost represents the cost breakdown for a chat
type Cost struct {
	Type             string      `json:"type"`
	Model            interface{} `json:"model"`
	PromptTokens     int         `json:"promptTokens"`
	CompletionTokens int         `json:"completionTokens"`
	Cost             float64     `json:"cost"`
}

// Assistant represents the assistant configuration
type Assistant struct {
	Transcriber                      *Transcriber                   `json:"transcriber,omitempty"`
	Model                            *Model                         `json:"model,omitempty"`
	Voice                            *Voice                         `json:"voice,omitempty"`
	FirstMessage                     *string                        `json:"firstMessage,omitempty"`
	FirstMessageInterruptionsEnabled *bool                          `json:"firstMessageInterruptionsEnabled,omitempty"`
	FirstMessageMode                 *string                        `json:"firstMessageMode,omitempty"`
	VoicemailDetection               *VoicemailDetection            `json:"voicemailDetection,omitempty"`
	ClientMessages                   *string                        `json:"clientMessages,omitempty"`
	ServerMessages                   *string                        `json:"serverMessages,omitempty"`
	MaxDurationSeconds               *int                           `json:"maxDurationSeconds,omitempty"`
	BackgroundSound                  *string                        `json:"backgroundSound,omitempty"`
	BackgroundDenoisingEnabled       *bool                          `json:"backgroundDenoisingEnabled,omitempty"`
	ModelOutputInMessagesEnabled     *bool                          `json:"modelOutputInMessagesEnabled,omitempty"`
	TransportConfigurations          []TransportConfiguration       `json:"transportConfigurations,omitempty"`
	ObservabilityPlan                *ObservabilityPlan             `json:"observabilityPlan,omitempty"`
	Credentials                      []Credential                   `json:"credentials,omitempty"`
	Hooks                            []Hook                         `json:"hooks,omitempty"`
	Name                             *string                        `json:"name,omitempty"`
	VoicemailMessage                 *string                        `json:"voicemailMessage,omitempty"`
	EndCallMessage                   *string                        `json:"endCallMessage,omitempty"`
	EndCallPhrases                   []string                       `json:"endCallPhrases,omitempty"`
	CompliancePlan                   *CompliancePlan                `json:"compliancePlan,omitempty"`
	Metadata                         map[string]interface{}         `json:"metadata,omitempty"`
	BackgroundSpeechDenoisingPlan    *BackgroundSpeechDenoisingPlan `json:"backgroundSpeechDenoisingPlan,omitempty"`
	AnalysisPlan                     *AnalysisPlan                  `json:"analysisPlan,omitempty"`
	ArtifactPlan                     *ArtifactPlan                  `json:"artifactPlan,omitempty"`
	MessagePlan                      *MessagePlan                   `json:"messagePlan,omitempty"`
	StartSpeakingPlan                *StartSpeakingPlan             `json:"startSpeakingPlan,omitempty"`
	StopSpeakingPlan                 *StopSpeakingPlan              `json:"stopSpeakingPlan,omitempty"`
	MonitorPlan                      *MonitorPlan                   `json:"monitorPlan,omitempty"`
	CredentialIDs                    []string                       `json:"credentialIds,omitempty"`
	Server                           *Server                        `json:"server,omitempty"`
	KeypadInputPlan                  *KeypadInputPlan               `json:"keypadInputPlan,omitempty"`
}

// AssistantOverrides represents overrides for assistant configuration
type AssistantOverrides struct {
	VariableValues map[string]interface{} `json:"variableValues,omitempty"`
	// Include all the same fields as Assistant for potential overrides
	Transcriber                      *Transcriber                   `json:"transcriber,omitempty"`
	Model                            *Model                         `json:"model,omitempty"`
	Voice                            *Voice                         `json:"voice,omitempty"`
	FirstMessage                     *string                        `json:"firstMessage,omitempty"`
	FirstMessageInterruptionsEnabled *bool                          `json:"firstMessageInterruptionsEnabled,omitempty"`
	FirstMessageMode                 *string                        `json:"firstMessageMode,omitempty"`
	VoicemailDetection               *VoicemailDetection            `json:"voicemailDetection,omitempty"`
	ClientMessages                   *string                        `json:"clientMessages,omitempty"`
	ServerMessages                   *string                        `json:"serverMessages,omitempty"`
	MaxDurationSeconds               *int                           `json:"maxDurationSeconds,omitempty"`
	BackgroundSound                  *string                        `json:"backgroundSound,omitempty"`
	BackgroundDenoisingEnabled       *bool                          `json:"backgroundDenoisingEnabled,omitempty"`
	ModelOutputInMessagesEnabled     *bool                          `json:"modelOutputInMessagesEnabled,omitempty"`
	TransportConfigurations          []TransportConfiguration       `json:"transportConfigurations,omitempty"`
	ObservabilityPlan                *ObservabilityPlan             `json:"observabilityPlan,omitempty"`
	Credentials                      []Credential                   `json:"credentials,omitempty"`
	Hooks                            []Hook                         `json:"hooks,omitempty"`
	Name                             *string                        `json:"name,omitempty"`
	VoicemailMessage                 *string                        `json:"voicemailMessage,omitempty"`
	EndCallMessage                   *string                        `json:"endCallMessage,omitempty"`
	EndCallPhrases                   []string                       `json:"endCallPhrases,omitempty"`
	CompliancePlan                   *CompliancePlan                `json:"compliancePlan,omitempty"`
	Metadata                         map[string]interface{}         `json:"metadata,omitempty"`
	BackgroundSpeechDenoisingPlan    *BackgroundSpeechDenoisingPlan `json:"backgroundSpeechDenoisingPlan,omitempty"`
	AnalysisPlan                     *AnalysisPlan                  `json:"analysisPlan,omitempty"`
	ArtifactPlan                     *ArtifactPlan                  `json:"artifactPlan,omitempty"`
	MessagePlan                      *MessagePlan                   `json:"messagePlan,omitempty"`
	StartSpeakingPlan                *StartSpeakingPlan             `json:"startSpeakingPlan,omitempty"`
	StopSpeakingPlan                 *StopSpeakingPlan              `json:"stopSpeakingPlan,omitempty"`
	MonitorPlan                      *MonitorPlan                   `json:"monitorPlan,omitempty"`
	CredentialIDs                    []string                       `json:"credentialIds,omitempty"`
	Server                           *Server                        `json:"server,omitempty"`
	KeypadInputPlan                  *KeypadInputPlan               `json:"keypadInputPlan,omitempty"`
}

// Transcriber represents transcriber configuration
type Transcriber struct {
	Provider                         string               `json:"provider"`
	Language                         *string              `json:"language,omitempty"`
	ConfidenceThreshold              *float64             `json:"confidenceThreshold,omitempty"`
	EnableUniversalStreamingAPI      *bool                `json:"enableUniversalStreamingApi,omitempty"`
	FormatTurns                      *bool                `json:"formatTurns,omitempty"`
	EndOfTurnConfidenceThreshold     *float64             `json:"endOfTurnConfidenceThreshold,omitempty"`
	MinEndOfTurnSilenceWhenConfident *int                 `json:"minEndOfTurnSilenceWhenConfident,omitempty"`
	WordFinalizationMaxWaitTime      *int                 `json:"wordFinalizationMaxWaitTime,omitempty"`
	MaxTurnSilence                   *int                 `json:"maxTurnSilence,omitempty"`
	RealtimeURL                      *string              `json:"realtimeUrl,omitempty"`
	WordBoost                        []string             `json:"wordBoost,omitempty"`
	EndUtteranceSilenceThreshold     *int                 `json:"endUtteranceSilenceThreshold,omitempty"`
	DisablePartialTranscripts        *bool                `json:"disablePartialTranscripts,omitempty"`
	FallbackPlan                     *TranscriberFallback `json:"fallbackPlan,omitempty"`
}

// TranscriberFallback represents fallback transcriber configuration
type TranscriberFallback struct {
	Transcribers []Transcriber `json:"transcribers"`
}

// Model represents the AI model configuration
type Model struct {
	Messages                  []ModelMessage  `json:"messages,omitempty"`
	Tools                     []Tool          `json:"tools,omitempty"`
	ToolIDs                   []string        `json:"toolIds,omitempty"`
	KnowledgeBase             *KnowledgeBase  `json:"knowledgeBase,omitempty"`
	KnowledgeBaseID           *string         `json:"knowledgeBaseId,omitempty"`
	Model                     string          `json:"model"`
	Provider                  string          `json:"provider"`
	Thinking                  *ThinkingConfig `json:"thinking,omitempty"`
	Temperature               *float64        `json:"temperature,omitempty"`
	MaxTokens                 *int            `json:"maxTokens,omitempty"`
	EmotionRecognitionEnabled *bool           `json:"emotionRecognitionEnabled,omitempty"`
	NumFastTurns              *int            `json:"numFastTurns,omitempty"`
}

// ModelMessage represents a message in the model context
type ModelMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

// Tool represents a tool configuration
type Tool struct {
	Messages               []ToolMessage           `json:"messages,omitempty"`
	Type                   string                  `json:"type"`
	Method                 *string                 `json:"method,omitempty"`
	TimeoutSeconds         *int                    `json:"timeoutSeconds,omitempty"`
	Name                   string                  `json:"name"`
	Description            *string                 `json:"description,omitempty"`
	URL                    *string                 `json:"url,omitempty"`
	Body                   *Schema                 `json:"body,omitempty"`
	Headers                *Schema                 `json:"headers,omitempty"`
	BackoffPlan            *BackoffPlan            `json:"backoffPlan,omitempty"`
	VariableExtractionPlan *VariableExtractionPlan `json:"variableExtractionPlan,omitempty"`
}

// ToolMessage represents a tool message
type ToolMessage struct {
	Contents   []MessageContent `json:"contents,omitempty"`
	Type       string           `json:"type"`
	Blocking   *bool            `json:"blocking,omitempty"`
	Content    *string          `json:"content,omitempty"`
	Conditions []Condition      `json:"conditions,omitempty"`
}

// MessageContent represents content within a message
type MessageContent struct {
	Type     string  `json:"type"`
	Text     *string `json:"text,omitempty"`
	Language *string `json:"language,omitempty"`
}

// Condition represents a condition for tool execution
type Condition struct {
	Operator string      `json:"operator"`
	Param    string      `json:"param"`
	Value    interface{} `json:"value"`
}

// Schema represents a JSON schema
type Schema struct {
	Type        *string                `json:"type,omitempty"`
	Items       map[string]interface{} `json:"items,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	Description *string                `json:"description,omitempty"`
	Pattern     *string                `json:"pattern,omitempty"`
	Format      *string                `json:"format,omitempty"`
	Required    []string               `json:"required,omitempty"`
	Enum        []string               `json:"enum,omitempty"`
	Title       *string                `json:"title,omitempty"`
}

// BackoffPlan represents retry configuration
type BackoffPlan struct {
	Type             string `json:"type"`
	MaxRetries       *int   `json:"maxRetries,omitempty"`
	BaseDelaySeconds *int   `json:"baseDelaySeconds,omitempty"`
}

// VariableExtractionPlan represents variable extraction configuration
type VariableExtractionPlan struct {
	Schema  *Schema         `json:"schema,omitempty"`
	Aliases []VariableAlias `json:"aliases,omitempty"`
}

// VariableAlias represents a variable alias
type VariableAlias struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// KnowledgeBase represents knowledge base configuration
type KnowledgeBase struct {
	Provider string  `json:"provider"`
	Server   *Server `json:"server,omitempty"`
}

// ThinkingConfig represents thinking configuration
type ThinkingConfig struct {
	Type         string `json:"type"`
	BudgetTokens *int   `json:"budgetTokens,omitempty"`
}

// Voice represents voice configuration
type Voice struct {
	CachingEnabled *bool          `json:"cachingEnabled,omitempty"`
	Provider       string         `json:"provider"`
	VoiceID        string         `json:"voiceId"`
	ChunkPlan      *ChunkPlan     `json:"chunkPlan,omitempty"`
	Speed          *float64       `json:"speed,omitempty"`
	FallbackPlan   *VoiceFallback `json:"fallbackPlan,omitempty"`
}

// ChunkPlan represents voice chunking configuration
type ChunkPlan struct {
	Enabled               *bool       `json:"enabled,omitempty"`
	MinCharacters         *int        `json:"minCharacters,omitempty"`
	PunctuationBoundaries *string     `json:"punctuationBoundaries,omitempty"`
	FormatPlan            *FormatPlan `json:"formatPlan,omitempty"`
}

// FormatPlan represents text formatting configuration
type FormatPlan struct {
	Enabled              *bool             `json:"enabled,omitempty"`
	NumberToDigitsCutoff *int              `json:"numberToDigitsCutoff,omitempty"`
	Replacements         []TextReplacement `json:"replacements,omitempty"`
	FormattersEnabled    *string           `json:"formattersEnabled,omitempty"`
}

// TextReplacement represents a text replacement rule
type TextReplacement struct {
	Type              string `json:"type"`
	ReplaceAllEnabled *bool  `json:"replaceAllEnabled,omitempty"`
	Key               string `json:"key"`
	Value             string `json:"value"`
}

// VoiceFallback represents voice fallback configuration
type VoiceFallback struct {
	Voices []Voice `json:"voices"`
}

// VoicemailDetection represents voicemail detection configuration
type VoicemailDetection struct {
	BeepMaxAwaitSeconds *int         `json:"beepMaxAwaitSeconds,omitempty"`
	Provider            string       `json:"provider"`
	BackoffPlan         *BackoffPlan `json:"backoffPlan,omitempty"`
}

// TransportConfiguration represents transport configuration
type TransportConfiguration struct {
	Provider          string  `json:"provider"`
	Timeout           *int    `json:"timeout,omitempty"`
	Record            *bool   `json:"record,omitempty"`
	RecordingChannels *string `json:"recordingChannels,omitempty"`
}

// ObservabilityPlan represents observability configuration
type ObservabilityPlan struct {
	Provider string                 `json:"provider"`
	Tags     []string               `json:"tags,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Credential represents API credentials
type Credential struct {
	Provider string  `json:"provider"`
	APIKey   string  `json:"apiKey"`
	Name     *string `json:"name,omitempty"`
}

// Hook represents a webhook configuration
type Hook struct {
	On      string       `json:"on"`
	Do      []HookAction `json:"do"`
	Filters []HookFilter `json:"filters,omitempty"`
}

// HookAction represents an action to perform in a hook
type HookAction struct {
	Type   string  `json:"type"`
	Tool   *Tool   `json:"tool,omitempty"`
	ToolID *string `json:"toolId,omitempty"`
}

// HookFilter represents a filter for hook execution
type HookFilter struct {
	Type  string   `json:"type"`
	Key   string   `json:"key"`
	OneOf []string `json:"oneOf,omitempty"`
}

// CompliancePlan represents compliance configuration
type CompliancePlan struct {
	HIPAAEnabled *HIPAAConfig `json:"hipaaEnabled,omitempty"`
	PCIEnabled   *PCIConfig   `json:"pciEnabled,omitempty"`
}

// HIPAAConfig represents HIPAA compliance configuration
type HIPAAConfig struct {
	HIPAAEnabled bool `json:"hipaaEnabled"`
}

// PCIConfig represents PCI compliance configuration
type PCIConfig struct {
	PCIEnabled bool `json:"pciEnabled"`
}

// BackgroundSpeechDenoisingPlan represents background speech denoising configuration
type BackgroundSpeechDenoisingPlan struct {
	SmartDenoisingPlan   *SmartDenoisingPlan   `json:"smartDenoisingPlan,omitempty"`
	FourierDenoisingPlan *FourierDenoisingPlan `json:"fourierDenoisingPlan,omitempty"`
}

// SmartDenoisingPlan represents smart denoising configuration
type SmartDenoisingPlan struct {
	Enabled bool `json:"enabled"`
}

// FourierDenoisingPlan represents Fourier denoising configuration
type FourierDenoisingPlan struct {
	Enabled               bool  `json:"enabled"`
	MediaDetectionEnabled *bool `json:"mediaDetectionEnabled,omitempty"`
	StaticThreshold       *int  `json:"staticThreshold,omitempty"`
	BaselineOffsetDB      *int  `json:"baselineOffsetDb,omitempty"`
	WindowSizeMS          *int  `json:"windowSizeMs,omitempty"`
	BaselinePercentile    *int  `json:"baselinePercentile,omitempty"`
}

// AnalysisPlan represents analysis configuration
type AnalysisPlan struct {
	MinMessagesThreshold    *int                      `json:"minMessagesThreshold,omitempty"`
	SummaryPlan             *SummaryPlan              `json:"summaryPlan,omitempty"`
	StructuredDataPlan      *StructuredDataPlan       `json:"structuredDataPlan,omitempty"`
	StructuredDataMultiPlan []StructuredDataMultiItem `json:"structuredDataMultiPlan,omitempty"`
	SuccessEvaluationPlan   *SuccessEvaluationPlan    `json:"successEvaluationPlan,omitempty"`
}

// SummaryPlan represents summary analysis configuration
type SummaryPlan struct {
	Messages       []interface{} `json:"messages,omitempty"`
	Enabled        *bool         `json:"enabled,omitempty"`
	TimeoutSeconds *int          `json:"timeoutSeconds,omitempty"`
}

// StructuredDataPlan represents structured data analysis configuration
type StructuredDataPlan struct {
	Messages       []interface{} `json:"messages,omitempty"`
	Enabled        *bool         `json:"enabled,omitempty"`
	Schema         *Schema       `json:"schema,omitempty"`
	TimeoutSeconds *int          `json:"timeoutSeconds,omitempty"`
}

// StructuredDataMultiItem represents a structured data multi-plan item
type StructuredDataMultiItem struct {
	Key  string              `json:"key"`
	Plan *StructuredDataPlan `json:"plan"`
}

// SuccessEvaluationPlan represents success evaluation configuration
type SuccessEvaluationPlan struct {
	Rubric         string        `json:"rubric"`
	Messages       []interface{} `json:"messages,omitempty"`
	Enabled        *bool         `json:"enabled,omitempty"`
	TimeoutSeconds *int          `json:"timeoutSeconds,omitempty"`
}

// ArtifactPlan represents artifact configuration
type ArtifactPlan struct {
	RecordingEnabled      *bool           `json:"recordingEnabled,omitempty"`
	RecordingFormat       *string         `json:"recordingFormat,omitempty"`
	VideoRecordingEnabled *bool           `json:"videoRecordingEnabled,omitempty"`
	PCAPEnabled           *bool           `json:"pcapEnabled,omitempty"`
	PCAPS3PathPrefix      *string         `json:"pcapS3PathPrefix,omitempty"`
	TranscriptPlan        *TranscriptPlan `json:"transcriptPlan,omitempty"`
	RecordingPath         *string         `json:"recordingPath,omitempty"`
}

// TranscriptPlan represents transcript configuration
type TranscriptPlan struct {
	Enabled       *bool   `json:"enabled,omitempty"`
	AssistantName *string `json:"assistantName,omitempty"`
	UserName      *string `json:"userName,omitempty"`
}

// MessagePlan represents message configuration
type MessagePlan struct {
	IdleMessages                             []string `json:"idleMessages,omitempty"`
	IdleMessageMaxSpokenCount                *int     `json:"idleMessageMaxSpokenCount,omitempty"`
	IdleMessageResetCountOnUserSpeechEnabled *bool    `json:"idleMessageResetCountOnUserSpeechEnabled,omitempty"`
	IdleTimeoutSeconds                       *int     `json:"idleTimeoutSeconds,omitempty"`
	SilenceTimeoutMessage                    *string  `json:"silenceTimeoutMessage,omitempty"`
}

// StartSpeakingPlan represents start speaking configuration
type StartSpeakingPlan struct {
	WaitSeconds                  *float64                      `json:"waitSeconds,omitempty"`
	SmartEndpointingEnabled      *bool                         `json:"smartEndpointingEnabled,omitempty"`
	SmartEndpointingPlan         *SmartEndpointingPlan         `json:"smartEndpointingPlan,omitempty"`
	CustomEndpointingRules       []CustomEndpointingRule       `json:"customEndpointingRules,omitempty"`
	TranscriptionEndpointingPlan *TranscriptionEndpointingPlan `json:"transcriptionEndpointingPlan,omitempty"`
}

// SmartEndpointingPlan represents smart endpointing configuration
type SmartEndpointingPlan struct {
	Provider string `json:"provider"`
}

// CustomEndpointingRule represents a custom endpointing rule
type CustomEndpointingRule struct {
	Type           string        `json:"type"`
	Regex          *string       `json:"regex,omitempty"`
	RegexOptions   []RegexOption `json:"regexOptions,omitempty"`
	TimeoutSeconds *int          `json:"timeoutSeconds,omitempty"`
}

// RegexOption represents a regex option
type RegexOption struct {
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}

// TranscriptionEndpointingPlan represents transcription endpointing configuration
type TranscriptionEndpointingPlan struct {
	OnPunctuationSeconds   *float64 `json:"onPunctuationSeconds,omitempty"`
	OnNoPunctuationSeconds *float64 `json:"onNoPunctuationSeconds,omitempty"`
	OnNumberSeconds        *float64 `json:"onNumberSeconds,omitempty"`
}

// StopSpeakingPlan represents stop speaking configuration
type StopSpeakingPlan struct {
	NumWords               *int     `json:"numWords,omitempty"`
	VoiceSeconds           *float64 `json:"voiceSeconds,omitempty"`
	BackoffSeconds         *float64 `json:"backoffSeconds,omitempty"`
	AcknowledgementPhrases []string `json:"acknowledgementPhrases,omitempty"`
	InterruptionPhrases    []string `json:"interruptionPhrases,omitempty"`
}

// MonitorPlan represents monitoring configuration
type MonitorPlan struct {
	ListenEnabled                *bool `json:"listenEnabled,omitempty"`
	ListenAuthenticationEnabled  *bool `json:"listenAuthenticationEnabled,omitempty"`
	ControlEnabled               *bool `json:"controlEnabled,omitempty"`
	ControlAuthenticationEnabled *bool `json:"controlAuthenticationEnabled,omitempty"`
}

// Server represents server configuration
type Server struct {
	TimeoutSeconds *int                   `json:"timeoutSeconds,omitempty"`
	URL            string                 `json:"url"`
	Headers        map[string]interface{} `json:"headers,omitempty"`
	BackoffPlan    *BackoffPlan           `json:"backoffPlan,omitempty"`
}

// KeypadInputPlan represents keypad input configuration
type KeypadInputPlan struct {
	Enabled        *bool   `json:"enabled,omitempty"`
	TimeoutSeconds *int    `json:"timeoutSeconds,omitempty"`
	Delimiters     *string `json:"delimiters,omitempty"`
}

// Session-related types for VAPI session management

// SessionResponse represents the response from creating or retrieving a session
type SessionResponse struct {
	ID          string `json:"id"`
	OrgID       string `json:"orgId"`
	AssistantID string `json:"assistantId"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
