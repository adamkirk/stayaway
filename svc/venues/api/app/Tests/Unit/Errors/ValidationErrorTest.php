<?php

namespace App\Tests\Errors;

use RuntimeException;
use Tests\UnitTestCase;
use App\Errors\ErrorType;
use App\Errors\ValidationError;

class ValidationErrorTest extends UnitTestCase
{
    public function test_setters()
    {
        $e = new ValidationError(
            'some_name',
            ErrorType::RecordNotFound,
            'some message',
        );

        $this->assertEquals('some_name', $e->name());
        $this->assertEquals(ErrorType::RecordNotFound, $e->type());
        $this->assertEquals('some message', $e->message());
    }

    /*
    Can't really test the default exception, it should be impossible to hit as 
    it's a backed enum.
    */
    public function test_message_defaults()
    {
        $notFound = new ValidationError('name', ErrorType::RecordNotFound);
        $this->assertEquals('record is not found', $notFound->message());

        $typeMismatch = new ValidationError('name', ErrorType::TypeMismatch);
        $this->assertEquals('the type for this value is incorrect', $typeMismatch->message());

        $valueNotAllowed = new ValidationError('name', ErrorType::ValueNotAllowed);
        $this->assertEquals('this value is not allowed', $valueNotAllowed->message());
    }
}
