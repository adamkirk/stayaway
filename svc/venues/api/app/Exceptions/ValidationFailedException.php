<?php

namespace App\Exceptions;

use Exception;
use Throwable;
use App\Collections\ValidationErrorCollection;

class ValidationFailedException extends Exception
{
    public function __construct(
        protected ValidationErrorCollection $errors,
        string $message = 'validation failed',
        int $code = 0,
        Throwable|null $previous = null
    ) {
        parent::__construct($message, $code, $previous);
    }

    public function errors(): ValidationErrorCollection
    {
        return $this->errors;
    }
}