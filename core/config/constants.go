package config

const (
	// AuthPrivateKey Token私匙
	AuthPrivateKey = `core_auth_private_key`
	// AuthPublicKey Token签名公钥
	AuthPublicKey = `core_auth_public_key`
	// AuthExpiry Token 过期时间
	AuthExpiry = `core_auth_expiry`

	// SmsTemplateCaptcha 验证码短信模板
	SmsTemplateCaptcha = `sys_sms_template_captcha`
	// EmailTemplateCaptcha 验证码邮件模板
	EmailTemplateCaptcha = `sys_email_template_captcha`
	// EmailSubjectCaptcha 验证码邮件主题
	EmailSubjectCaptcha = `sys_email_subject_captcha`

	// SMTPHost 邮件服务器地址
	SMTPHost = `sys_smtp_host`
	// SMTPPort 邮件服务器端口
	SMTPPort = `sys_smtp_port`
	// SMTPUsername 系统邮件发送账号
	SMTPUsername = `sys_smpt_username`
	// SMTPPassword 系统邮件发送密码
	SMTPPassword = `sys_smpt_password`

	// MediaURI 媒体路径
	MediaURI = `core_media_uri`
	// CodecPrefs 媒体编码
	CodecPrefs = `core_codec_prefs`
	// MediaMixInboundOutboundCodecs 协商全部媒体编码
	MediaMixInboundOutboundCodecs = `core_media_mix_inbound_outbound_codecs`
	// RecordingPath 录音路径
	RecordingPath = `core_recording_path`
	// RecordingType 录音转换类型
	RecordingType = `core_recording_type`

	// CallTracking 是否开启呼叫跟踪
	CallTracking = `core_call_tracking`

	// AgentStateInboundCall 是否启用坐席呼入状态跟踪
	AgentStateInboundCall = `agent_state_inbound_call`

	// CronExpireDays 定时数据保存时间
	CronExpireDays = `cron_expire_days`
	// PrimaryDomainName 第三方域
	PrimaryDomainName = `core_primary_domain_name`

	// UniqueEmployeeNumber 工号是否全局唯一
	UniqueEmployeeNumber = `core_unique_employee_number`

	// ServiceNames 服务名字列表
	ServiceNames = `sys_service_name`

	// CacheSettings 缓存配置
	CacheSettings = `core_cache_settings`

	// AgentFailAcwTime 坐席桥接失败后话后处理时长
	AgentFailAcwTime = `agent_fail_acw_time`

	//是否开启通话计时挂断：true：开启，false：关闭
	CallTimeStart = `core_call_time_start`
)
