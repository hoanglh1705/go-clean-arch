package commonhelper

type ContentKeyType string

const (
	ContextKeyType_TraceId    ContentKeyType = "TraceId"
	ContextKeyType_TimeMs     ContentKeyType = "TimeMs"
	HeaderKeyType_UserId      ContentKeyType = "UId"
	HeaderKeyType_FeUserId    ContentKeyType = "FeUId"
	ContextKeyType_Subject    ContentKeyType = "Sub"
	ContextKeyType_AppSubject ContentKeyType = "AppSub"
	HeaderKeyType_Token       ContentKeyType = "Token"
)

type MetadataKeyType string

const (
	MetadataKeyType_Authorization MetadataKeyType = "authorization"
	MetadataKeyType_Subject       MetadataKeyType = "sub"
	MetadataKeyType_AppSubject    MetadataKeyType = "appsub"
)

type TimezoneType string

const (
	TimezoneType_AsiaHoChiMinh TimezoneType = "Asia/Ho_Chi_Minh"
)

type PayloadTrackingField string

const (
	PayloadTrackingField_TraceId PayloadTrackingField = "traceId"
)
