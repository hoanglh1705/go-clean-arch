package commonhelper

type NotiChannelType string

const (
	NotiChannelType_SMS              NotiChannelType = "sms"
	NotiChannelType_Email            NotiChannelType = "email"
	NotiChannelType_PushNotification NotiChannelType = "push_notification"
	NotiChannelType_Webhook          NotiChannelType = "webhook"
)

type NotiEventCodeType string

const (
	NotiEventCodeType_LinkSlackAccountOtp  NotiEventCodeType = "LinkSlackAccountOtp"
	NotiEventCodeType_StockNotiBigVol      NotiEventCodeType = "STOCK_NOTIFICATIONS_BIG_VOL"
	NotiEventCodeType_StockNotiPriceChange NotiEventCodeType = "STOCK_NOTIFICATIONS_PRICE_CHANGE"
)

type NotiTopicType string

const (
	NotiTopicType_StockNotifications NotiTopicType = "stock_notifications"
)

type NotiDeliveryPriorityType string

const (
	NotiDeliveryPriorityType_High   NotiDeliveryPriorityType = "HIGH"
	NotiDeliveryPriorityType_Normal NotiDeliveryPriorityType = "NORMAL"
	NotiDeliveryPriorityType_Low    NotiDeliveryPriorityType = "LOW"
)

type NotiDeliveryTopicType string

const (
	NotiDeliveryTopicTestSMS              NotiDeliveryTopicType = "test_sms"
	NotiDeliveryTopicTestEmail            NotiDeliveryTopicType = "test_email"
	NotiDeliveryTopicTestPushNotification NotiDeliveryTopicType = "test_push_notification"
	NotiDeliveryTopicTestWebhook          NotiDeliveryTopicType = "test_webhook"
	NotiDeliveryTopicGeneralNotification  NotiDeliveryTopicType = "general_notifications"
	NotiDeliveryTopicStockNotification    NotiDeliveryTopicType = "stock_notifications"
)

type NotiDeliveryStatusType string

const (
	NotiDeliveryStatusType_New    NotiDeliveryStatusType = "NEW"
	NotiDeliveryStatusType_Sent   NotiDeliveryStatusType = "SENT"
	NotiDeliveryStatusType_Failed NotiDeliveryStatusType = "FAILED"
)
