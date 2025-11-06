package constants

type BaseResponse struct {
	Code    int
	Message string
	Data    interface{}
}

type StatusCode interface {
	GetCode() int
	GetMessage() string
}

type UserStatusCode int

const (
	UserSuccess      UserStatusCode = 200
	UserCreated      UserStatusCode = 201
	UserBadRequest   UserStatusCode = 400
	UserUnauthorized UserStatusCode = 401
	UserForbidden    UserStatusCode = 403
	UserNotFound     UserStatusCode = 404
	UserConflict     UserStatusCode = 409
	UserSystemError  UserStatusCode = 500
)

func (u UserStatusCode) GetCode() int {
	return int(u)
}

func (u UserStatusCode) GetMessage() string {
	switch u {
	case UserSuccess:
		return "操作成功"
	case UserCreated:
		return "创建成功"
	case UserBadRequest:
		return "参数错误"
	case UserUnauthorized:
		return "未授权"
	case UserForbidden:
		return "禁止访问"
	case UserNotFound:
		return "用户不存在"
	case UserConflict:
		return "用户已存在"
	case UserSystemError:
		return "系统错误"
	default:
		return "未知错误"
	}
}

// 通用状态码
type CommonStatusCode int

const (
	Success      CommonStatusCode = 200
	Created      CommonStatusCode = 201
	BadRequest   CommonStatusCode = 400
	Unauthorized CommonStatusCode = 401
	Forbidden    CommonStatusCode = 403
	NotFound     CommonStatusCode = 404
	Conflict     CommonStatusCode = 409
	SystemError  CommonStatusCode = 500
)

func (c CommonStatusCode) GetCode() int { return int(c) }

func (c CommonStatusCode) GetMessage() string {
	switch c {
	case Success:
		return "success"
	case Created:
		return "created"
	case BadRequest:
		return "bad request"
	case Unauthorized:
		return "unauthorized"
	case Forbidden:
		return "forbidden"
	case NotFound:
		return "not found"
	case Conflict:
		return "conflict"
	case SystemError:
		return "system error"
	default:
		return "unknown status"
	}
}

func BuildResponseWithStatus(status StatusCode, data interface{}) BaseResponse {
	return BaseResponse{
		Code:    status.GetCode(),
		Message: status.GetMessage(),
		Data:    data,
	}
}

// ContentStatusCode 内容相关状态码
type ContentStatusCode int

const (
	ContentSuccess      ContentStatusCode = 200
	ContentCreated      ContentStatusCode = 201
	ContentBadRequest   ContentStatusCode = 400
	ContentUnauthorized ContentStatusCode = 401
	ContentForbidden    ContentStatusCode = 403
	ContentNotFound     ContentStatusCode = 404
	ContentConflict     ContentStatusCode = 409
	ContentSystemError  ContentStatusCode = 500
)

func (c ContentStatusCode) GetCode() int { return int(c) }

func (c ContentStatusCode) GetMessage() string {
	switch c {
	case ContentSuccess:
		return "操作成功"
	case ContentCreated:
		return "创建成功"
	case ContentBadRequest:
		return "参数错误"
	case ContentUnauthorized:
		return "未授权"
	case ContentForbidden:
		return "禁止访问"
	case ContentNotFound:
		return "内容不存在"
	case ContentConflict:
		return "冲突错误"
	case ContentSystemError:
		return "系统错误"
	default:
		return "未知错误"
	}
}

// CommentStatusCode 评论相关状态码
type CommentStatusCode int

const (
	CommentSuccess      CommentStatusCode = 200
	CommentCreated      CommentStatusCode = 201
	CommentBadRequest   CommentStatusCode = 400
	CommentUnauthorized CommentStatusCode = 401
	CommentForbidden    CommentStatusCode = 403
	CommentNotFound     CommentStatusCode = 404
	CommentConflict     CommentStatusCode = 409
	CommentSystemError  CommentStatusCode = 500
)

func (c CommentStatusCode) GetCode() int { return int(c) }

func (c CommentStatusCode) GetMessage() string {
	switch c {
	case CommentSuccess:
		return "操作成功"
	case CommentCreated:
		return "创建成功"
	case CommentBadRequest:
		return "参数错误"
	case CommentUnauthorized:
		return "未授权"
	case CommentForbidden:
		return "禁止访问"
	case CommentNotFound:
		return "评论不存在"
	case CommentConflict:
		return "冲突错误"
	case CommentSystemError:
		return "系统错误"
	default:
		return "未知错误"
	}
}

// TagStatusCode 标签相关状态码
type TagStatusCode int

const (
	TagSuccess      TagStatusCode = 200
	TagCreated      TagStatusCode = 201
	TagBadRequest   TagStatusCode = 400
	TagUnauthorized TagStatusCode = 401
	TagForbidden    TagStatusCode = 403
	TagNotFound     TagStatusCode = 404
	TagConflict     TagStatusCode = 409
	TagSystemError  TagStatusCode = 500
)

func (t TagStatusCode) GetCode() int { return int(t) }

func (t TagStatusCode) GetMessage() string {
	switch t {
	case TagSuccess:
		return "操作成功"
	case TagCreated:
		return "创建成功"
	case TagBadRequest:
		return "参数错误"
	case TagUnauthorized:
		return "未授权"
	case TagForbidden:
		return "禁止访问"
	case TagNotFound:
		return "标签不存在"
	case TagConflict:
		return "冲突错误"
	case TagSystemError:
		return "系统错误"
	default:
		return "未知错误"
	}
}

// EmailStatusCode 邮件相关状态码
type EmailStatusCode int

const (
	EmailSuccess      EmailStatusCode = 200
	EmailCreated      EmailStatusCode = 201
	EmailBadRequest   EmailStatusCode = 400
	EmailUnauthorized EmailStatusCode = 401
	EmailForbidden    EmailStatusCode = 403
	EmailNotFound     EmailStatusCode = 404
	EmailConflict     EmailStatusCode = 409
	EmailSystemError  EmailStatusCode = 500
)

func (e EmailStatusCode) GetCode() int { return int(e) }

func (e EmailStatusCode) GetMessage() string {
	switch e {
	case EmailSuccess:
		return "操作成功"
	case EmailCreated:
		return "发送成功"
	case EmailBadRequest:
		return "参数错误"
	case EmailUnauthorized:
		return "未授权"
	case EmailForbidden:
		return "禁止访问"
	case EmailNotFound:
		return "邮件记录不存在"
	case EmailConflict:
		return "冲突错误"
	case EmailSystemError:
		return "系统错误"
	default:
		return "未知错误"
	}
}

type FileStatusCode int

const (
	FFileUploadSuccess FileStatusCode = 200
	FileUploadFailed   FileStatusCode = 410
)

func (f FileStatusCode) GetCode() int { return int(f) }

func (f FileStatusCode) GetMessage() string {
	switch f {
	case FFileUploadSuccess:
		return "文件上传成功"
	case FileUploadFailed:
		return "文件上传失败"
	default:
		return "未知错误"
	}
}
