package style

var BoldText = New(Bold)
var Header = New(Bold, Underline)

var Success = New(Green)
var BoldSuccess = New(Bold, Green)

var Error = New(Red)
var BoldError = New(Bold, Red)

var File = New(Yellow)
var BoldFile = New(Bold, Yellow)

var FileV2 = New(Magenta)
var BoldFileV2 = New(Bold, Magenta)

var Create = Success
var BoldCreate = BoldSuccess

var Delete = Error
var BoldDelete = BoldError

var Info = New(Blue)
var BoldInfo = New(Bold, Blue)

var InfoV2 = New(Cyan)
var BoldInfoV2 = New(Bold, Cyan)
