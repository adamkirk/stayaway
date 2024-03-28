<?php

namespace App\Exceptions;

use Throwable;
use \Exception;

class InvalidPropertyException extends Exception
{
    public function __construct(string $object, string $field, string $message, int $code = 0, Throwable|null $previous = null)
    {
        parent::__construct(
            "Invalid value for '$field' on '$object': $message",
            $code, 
            $previous,
        );
    }
}