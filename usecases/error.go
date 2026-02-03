package usecases

import "errors"

// general error
var (
	ErrInternalServerError = errors.New("เกิดข้อผิดพลาดภายในระบบ")
	ErrInvalidInput        = errors.New("ข้อมูลนำเข้าไม่ถูกต้อง")
	ErrInvalidID           = errors.New("รูปแบบ ID ไม่ถูกต้อง")
	ErrInvalidDate         = errors.New("รูปแบบวันที่ไม่ถูกต้อง (ต้องเป็น YYYY-MM-DD)")
	ErrCantFuture          = errors.New("วันที่ไม่สามารถเป็นอนาคตได้")
)

// auth, user
var (
	ErrUserNotFound         = errors.New("user not found")
	ErrCurrentPasswordWrong = errors.New("current password incorrect")
	ErrPasswordHashFailed   = errors.New("failed to hash password")
	ErrApproved             = errors.New("email is already verify")
	ErrFailedToFindEmail    = errors.New("failed to find email")
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrEmailNotVerify       = errors.New("email Not Verify")
	ErrAccountPending       = errors.New("account pending")
	ErrAccountReject        = errors.New("your account was rejected")
	ErrFailed               = errors.New("cannot save user")
	ErrFailedToSendEmail    = errors.New("failed to send email")
	ErrEmailAlready         = errors.New("email is already")
	ErrInvalidOrExpiry      = errors.New("token invalid or expired")
)

// pig
var (
	ErrPigNotFound          = errors.New("ไม่เจอหมูรหัสนี้ในระบบ")
	ErrPigCodeAlreadyExists = errors.New("รหัสหมูนี้มีอยู่ในระบบแล้ว")
	ErrInvalidWeight        = errors.New("น้ำหนักต้องมากกว่า 0")
	ErrMaleAsMother         = errors.New("เพศผู้ไม่สามารถเป็นเเม่พันธุ์ได้")
	ErrFemaleAsFather       = errors.New("เพศเมียไม่สามารถเป็นพ่อพันธุ์ได้")
	ErrMaleInvalidStatus    = errors.New("เพศผู้ไม่สามารถมีสถานะอุ้มท้องหรือให้นมลูกได้")
	ErrPigTypeInvalidStatus = errors.New("ประเภทหมูไม่สามารถมีสถานะนี้ได้")
	ErrFailedUpdate         = errors.New("ไม่สามารถอัปเดตได้")
	ErrIsUsedInBreeding     = errors.New("ไม่สามารถลบได้ หมูถูกใช้ในการผสมพันธุ์เเล้ว")
)

// breeding
var (
	ErrBreedingNotFound = errors.New("ไม่พบข้อมูลการผสมพันธุ์")
	ErrFatherNotFound   = errors.New("ไม่พบข้อมูลพ่อพันธุ์")
	ErrMotherNotFound   = errors.New("ไม่พบข้อมูลแม่พันธุ์")

	// Validation Logic
	ErrSamePig              = errors.New("พ่อพันธุ์และแม่พันธุ์ต้องไม่ใช่ตัวเดียวกัน")
	ErrInvalidFatherBreeder = errors.New("พ่อพันธุ์ต้องเป็นเพศผู้และเป็นประเภทพ่อพันธุ์เท่านั้น")
	ErrInvalidMotherBreeder = errors.New("แม่พันธุ์ต้องเป็นเพศเมียและเป็นประเภทแม่พันธุ์เท่านั้น")
	ErrPigNotReady          = errors.New("หมูยังไม่พร้อมผสมพันธุ์ (สถานะต้องเป็น 'พร้อมผสม')")
	ErrDuplicateBreeding    = errors.New("มีการบันทึกการผสมพันธุ์ของคู่นี้ในวันนี้ไปแล้ว")
)

// feeding
var (
	ErrFeedingNotFound = errors.New("ไม่พบข้อมูลการให้อาหาร")
	ErrFoodNotZero     = errors.New("ปริมาณอาหารต้องมากกว่า 0")
	ErrFoodNotFound    = errors.New("ไม่พบข้อมูลอาหารในคลัง")
	ErrNotEnoughFood   = errors.New("ปริมาณอาหารในคลังไม่เพียงพอ")
	ErrNoValidPigs     = errors.New("ไม่มีหมูที่สามารถให้อาหารได้ในรายการที่เลือก (หมูอาจจะขายหรือตายไปแล้ว)")
)

// stock
var (
	ErrFoodStockNotFound      = errors.New("ไม่พบข้อมูลสต็อกอาหารที่ระบุ")
	ErrFoodStockAlreadyExists = errors.New("มีสต็อกของอาหารชนิดนี้อยู่แล้ว")
	ErrNotEnoughFoodStock     = errors.New("ปริมาณอาหารในสต็อกไม่เพียงพอ")
	ErrInvalidStockAmount     = errors.New("ปริมาณอาหารต้องมากกว่า 0 และไม่ติดลบ")
	ErrFoodStockUsed          = errors.New("ไม่สามารถลบรายการนี้ได้ เนื่องจากมีการนำไปใช้งานในการให้อาหารแล้ว")
)
