<?php

namespace App\Errors;

use RuntimeException;
use App\Errors\ErrorType;

class ValidationError
{
    public function __construct(
        protected string $name,
        protected ErrorType $errorType,
        protected string|null $message = null,
    ) {}

    public function name(): string
    {
        return $this->name;
    }

    public function type(): ErrorType
    {
        return $this->errorType;
    }

    public function message(): string
    {
        if ($this->message !== null) {
            return $this->message;
        }

        return match($this->errorType) {
            ErrorType::RecordNotFound => "record is not found",
            ErrorType::TypeMismatch => "the type for this value is incorrect",
            ErrorType::ValueNotAllowed => "this value is not allowed",
            default => throw new RuntimeException("Failed to build default error message for error type: " . $this->errorType->value),
        };
    }
}