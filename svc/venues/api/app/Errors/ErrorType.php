<?php

namespace App\Errors;

enum ErrorType: string {
    case TypeMismatch = 'type_mismatch';
    case ValueNotAllowed = 'value_not_allowed';
    case RecordNotFound = 'record_not_found';
    case Required = 'required';
}