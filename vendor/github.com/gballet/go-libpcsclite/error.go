// BSD 3-Clause License
//
// Copyright (c) 2019, Guillaume Ballet
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of the copyright holder nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package pcsc

import "fmt"

type ErrorCode uint32

const (
	SCardSuccess                   ErrorCode = 0x00000000 /* No error was encountered. */
	ErrSCardInternal               ErrorCode = 0x80100001 /* An internal consistency check failed. */
	ErrSCardCancelled              ErrorCode = 0x80100002 /* The action was cancelled by an SCardCancel request. */
	ErrSCardInvalidHandle          ErrorCode = 0x80100003 /* The supplied handle was invalid. */
	ErrSCardInvalidParameter       ErrorCode = 0x80100004 /* One or more of the supplied parameters could not be properly interpreted. */
	ErrSCardInvalidTarget          ErrorCode = 0x80100005 /* Registry startup information is missing or invalid. */
	ErrSCardNoMemory               ErrorCode = 0x80100006 /* Not enough memory available to complete this command. */
	ErrSCardWaitedTooLong          ErrorCode = 0x80100007 /* An internal consistency timer has expired. */
	ErrSCardInsufficientBuffer     ErrorCode = 0x80100008 /* The data buffer to receive returned data is too small for the returned data. */
	ErrScardUnknownReader          ErrorCode = 0x80100009 /* The specified reader name is not recognized. */
	ErrSCardTimeout                ErrorCode = 0x8010000A /* The user-specified timeout value has expired. */
	ErrSCardSharingViolation       ErrorCode = 0x8010000B /* The smart card cannot be accessed because of other connections outstanding. */
	ErrSCardNoSmartCard            ErrorCode = 0x8010000C /* The operation requires a Smart Card, but no Smart Card is currently in the device. */
	ErrSCardUnknownCard            ErrorCode = 0x8010000D /* The specified smart card name is not recognized. */
	ErrSCardCannotDispose          ErrorCode = 0x8010000E /* The system could not dispose of the media in the requested manner. */
	ErrSCardProtoMismatch          ErrorCode = 0x8010000F /* The requested protocols are incompatible with the protocol currently in use with the smart card. */
	ErrSCardNotReady               ErrorCode = 0x80100010 /* The reader or smart card is not ready to accept commands. */
	ErrSCardInvalidValue           ErrorCode = 0x80100011 /* One or more of the supplied parameters values could not be properly interpreted. */
	ErrSCardSystemCancelled        ErrorCode = 0x80100012 /* The action was cancelled by the system, presumably to log off or shut down. */
	ErrSCardCommError              ErrorCode = 0x80100013 /* An internal communications error has been detected. */
	ErrScardUnknownError           ErrorCode = 0x80100014 /* An internal error has been detected, but the source is unknown. */
	ErrSCardInvalidATR             ErrorCode = 0x80100015 /* An ATR obtained from the registry is not a valid ATR string. */
	ErrSCardNotTransacted          ErrorCode = 0x80100016 /* An attempt was made to end a non-existent transaction. */
	ErrSCardReaderUnavailable      ErrorCode = 0x80100017 /* The specified reader is not currently available for use. */
	ErrSCardShutdown               ErrorCode = 0x80100018 /* The operation has been aborted to allow the server application to exit. */
	ErrSCardPCITooSmall            ErrorCode = 0x80100019 /* The PCI Receive buffer was too small. */
	ErrSCardReaderUnsupported      ErrorCode = 0x8010001A /* The reader driver does not meet minimal requirements for support. */
	ErrSCardDuplicateReader        ErrorCode = 0x8010001B /* The reader driver did not produce a unique reader name. */
	ErrSCardCardUnsupported        ErrorCode = 0x8010001C /* The smart card does not meet minimal requirements for support. */
	ErrScardNoService              ErrorCode = 0x8010001D /* The Smart card resource manager is not running. */
	ErrSCardServiceStopped         ErrorCode = 0x8010001E /* The Smart card resource manager has shut down. */
	ErrSCardUnexpected             ErrorCode = 0x8010001F /* An unexpected card error has occurred. */
	ErrSCardUnsupportedFeature     ErrorCode = 0x8010001F /* This smart card does not support the requested feature. */
	ErrSCardICCInstallation        ErrorCode = 0x80100020 /* No primary provider can be found for the smart card. */
	ErrSCardICCCreateOrder         ErrorCode = 0x80100021 /* The requested order of object creation is not supported. */
	ErrSCardDirNotFound            ErrorCode = 0x80100023 /* The identified directory does not exist in the smart card. */
	ErrSCardFileNotFound           ErrorCode = 0x80100024 /* The identified file does not exist in the smart card. */
	ErrSCardNoDir                  ErrorCode = 0x80100025 /* The supplied path does not represent a smart card directory. */
	ErrSCardNoFile                 ErrorCode = 0x80100026 /* The supplied path does not represent a smart card file. */
	ErrScardNoAccess               ErrorCode = 0x80100027 /* Access is denied to this file. */
	ErrSCardWriteTooMany           ErrorCode = 0x80100028 /* The smart card does not have enough memory to store the information. */
	ErrSCardBadSeek                ErrorCode = 0x80100029 /* There was an error trying to set the smart card file object pointer. */
	ErrSCardInvalidCHV             ErrorCode = 0x8010002A /* The supplied PIN is incorrect. */
	ErrSCardUnknownResMNG          ErrorCode = 0x8010002B /* An unrecognized error code was returned from a layered component. */
	ErrSCardNoSuchCertificate      ErrorCode = 0x8010002C /* The requested certificate does not exist. */
	ErrSCardCertificateUnavailable ErrorCode = 0x8010002D /* The requested certificate could not be obtained. */
	ErrSCardNoReadersAvailable     ErrorCode = 0x8010002E /* Cannot find a smart card reader. */
	ErrSCardCommDataLost           ErrorCode = 0x8010002F /* A communications error with the smart card has been detected. Retry the operation. */
	ErrScardNoKeyContainer         ErrorCode = 0x80100030 /* The requested key container does not exist on the smart card. */
	ErrSCardServerTooBusy          ErrorCode = 0x80100031 /* The Smart Card Resource Manager is too busy to complete this operation. */
	ErrSCardUnsupportedCard        ErrorCode = 0x80100065 /* The reader cannot communicate with the card, due to ATR string configuration conflicts. */
	ErrSCardUnresponsiveCard       ErrorCode = 0x80100066 /* The smart card is not responding to a reset. */
	ErrSCardUnpoweredCard          ErrorCode = 0x80100067 /* Power has been removed from the smart card, so that further communication is not possible. */
	ErrSCardResetCard              ErrorCode = 0x80100068 /* The smart card has been reset, so any shared state information is invalid. */
	ErrSCardRemovedCard            ErrorCode = 0x80100069 /* The smart card has been removed, so further communication is not possible. */
	ErrSCardSecurityViolation      ErrorCode = 0x8010006A /* Access was denied because of a security violation. */
	ErrSCardWrongCHV               ErrorCode = 0x8010006B /* The card cannot be accessed because the wrong PIN was presented. */
	ErrSCardCHVBlocked             ErrorCode = 0x8010006C /* The card cannot be accessed because the maximum number of PIN entry attempts has been reached. */
	ErrSCardEOF                    ErrorCode = 0x8010006D /* The end of the smart card file has been reached. */
	ErrSCardCancelledByUser        ErrorCode = 0x8010006E /* The user pressed "Cancel" on a Smart Card Selection Dialog. */
	ErrSCardCardNotAuthenticated   ErrorCode = 0x8010006F /* No PIN was presented to the smart card. */
)

// Code returns the error code, with an uint32 type to be used in PutUInt32
func (code ErrorCode) Code() uint32 {
	return uint32(code)
}

func (code ErrorCode) Error() error {
	switch code {
	case SCardSuccess:
		return fmt.Errorf("command successful")

	case ErrSCardInternal:
		return fmt.Errorf("internal error")

	case ErrSCardCancelled:
		return fmt.Errorf("command cancelled")

	case ErrSCardInvalidHandle:
		return fmt.Errorf("invalid handle")

	case ErrSCardInvalidParameter:
		return fmt.Errorf("invalid parameter given")

	case ErrSCardInvalidTarget:
		return fmt.Errorf("invalid target given")

	case ErrSCardNoMemory:
		return fmt.Errorf("not enough memory")

	case ErrSCardWaitedTooLong:
		return fmt.Errorf("waited too long")

	case ErrSCardInsufficientBuffer:
		return fmt.Errorf("insufficient buffer")

	case ErrScardUnknownReader:
		return fmt.Errorf("unknown reader specified")

	case ErrSCardTimeout:
		return fmt.Errorf("command timeout")

	case ErrSCardSharingViolation:
		return fmt.Errorf("sharing violation")

	case ErrSCardNoSmartCard:
		return fmt.Errorf("no smart card inserted")

	case ErrSCardUnknownCard:
		return fmt.Errorf("unknown card")

	case ErrSCardCannotDispose:
		return fmt.Errorf("cannot dispose handle")

	case ErrSCardProtoMismatch:
		return fmt.Errorf("card protocol mismatch")

	case ErrSCardNotReady:
		return fmt.Errorf("subsystem not ready")

	case ErrSCardInvalidValue:
		return fmt.Errorf("invalid value given")

	case ErrSCardSystemCancelled:
		return fmt.Errorf("system cancelled")

	case ErrSCardCommError:
		return fmt.Errorf("rpc transport error")

	case ErrScardUnknownError:
		return fmt.Errorf("unknown error")

	case ErrSCardInvalidATR:
		return fmt.Errorf("invalid ATR")

	case ErrSCardNotTransacted:
		return fmt.Errorf("transaction failed")

	case ErrSCardReaderUnavailable:
		return fmt.Errorf("reader is unavailable")

	/* case SCARD_P_SHUTDOWN: */
	case ErrSCardPCITooSmall:
		return fmt.Errorf("PCI struct too small")

	case ErrSCardReaderUnsupported:
		return fmt.Errorf("reader is unsupported")

	case ErrSCardDuplicateReader:
		return fmt.Errorf("reader already exists")

	case ErrSCardCardUnsupported:
		return fmt.Errorf("card is unsupported")

	case ErrScardNoService:
		return fmt.Errorf("service not available")

	case ErrSCardServiceStopped:
		return fmt.Errorf("service was stopped")

	/* case SCARD_E_UNEXPECTED: */
	/* case SCARD_E_ICC_CREATEORDER: */
	/* case SCARD_E_UNSUPPORTED_FEATURE: */
	/* case SCARD_E_DIR_NOT_FOUND: */
	/* case SCARD_E_NO_DIR: */
	/* case SCARD_E_NO_FILE: */
	/* case SCARD_E_NO_ACCESS: */
	/* case SCARD_E_WRITE_TOO_MANY: */
	/* case SCARD_E_BAD_SEEK: */
	/* case SCARD_E_INVALID_CHV: */
	/* case SCARD_E_UNKNOWN_RES_MNG: */
	/* case SCARD_E_NO_SUCH_CERTIFICATE: */
	/* case SCARD_E_CERTIFICATE_UNAVAILABLE: */
	case ErrSCardNoReadersAvailable:
		return fmt.Errorf("cannot find a smart card reader")

	/* case SCARD_E_COMM_DATA_LOST: */
	/* case SCARD_E_NO_KEY_CONTAINER: */
	/* case SCARD_E_SERVER_TOO_BUSY: */
	case ErrSCardUnsupportedCard:
		return fmt.Errorf("Card is not supported")

	case ErrSCardUnresponsiveCard:
		return fmt.Errorf("Card is unresponsive")

	case ErrSCardUnpoweredCard:
		return fmt.Errorf("Card is unpowered")

	case ErrSCardResetCard:
		return fmt.Errorf("Card was reset")

	case ErrSCardRemovedCard:
		return fmt.Errorf("Card was removed")

	/* case SCARD_W_SECURITY_VIOLATION: */
	/* case SCARD_W_WRONG_CHV: */
	/* case SCARD_W_CHV_BLOCKED: */
	/* case SCARD_W_EOF: */
	/* case SCARD_W_CANCELLED_BY_USER: */
	/* case SCARD_W_CARD_NOT_AUTHENTICATED: */

	case ErrSCardUnsupportedFeature:
		return fmt.Errorf("feature not supported")

	default:
		return fmt.Errorf("unknown error: %08x", code)
	}
}
