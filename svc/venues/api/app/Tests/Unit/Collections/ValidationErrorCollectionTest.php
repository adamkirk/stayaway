<?php

namespace App\Tests\Unit\Collections;

use Tests\UnitTestCase;
use App\Errors\ErrorType;
use App\Errors\ValidationError;
use App\Collections\ValidationErrorCollection;

class ValidationErrorCollectionTest extends UnitTestCase
{

    protected function mockValidationError(string $name, ErrorType $errorType, ?string $message): ValidationError
    {
        $mock = $this->getMockBuilder(ValidationError::class)
            ->disableOriginalConstructor()
            ->getMock();

        $mock->expects($this->any())->method('name')->willReturn($name);
        $mock->expects($this->any())->method('type')->willReturn($errorType);
        $mock->expects($this->any())->method('message')->willReturn($message ?? 'default message');

        return $mock;
    }

    public function test_errors_can_be_added()
    {
        $coll = ValidationErrorCollection::new();

        $error1 = $this->mockValidationError(
            'field_1',
            ErrorType::ValueNotAllowed,
            "some message",
        );
        $error2 = $this->mockValidationError(
            'field_1',
            ErrorType::ValueNotAllowed,
            null
        );

        $coll->add($error1);
        $coll->add($error2);

        $errors = $coll->all();

        $this->assertCount(2, $errors);
        $this->assertSame($error1, $errors[0]);
        $this->assertSame($error2, $errors[1]);
    }

    public function test_conversion_to_array()
    {
        $error1 = $this->mockValidationError(
            'field_1',
            ErrorType::ValueNotAllowed,
            "some message",
        );
        $error2 = $this->mockValidationError(
            'field_1',
            ErrorType::ValueNotAllowed,
            "some other message",
        );
        $error3 = $this->mockValidationError(
            'field_2',
            ErrorType::ValueNotAllowed,
            null
        );
        $error4 = $this->mockValidationError(
            'field_3.sub',
            ErrorType::ValueNotAllowed,
            "sub error"
        );

        $coll = ValidationErrorCollection::new(
            $error1,
            $error2,
            $error3,
            $error4,
        );

        $this->assertEquals([
            'field_1' => [
                'some message',
                'some other message',
            ],
            'field_2' => [
                'default message',
            ],
            'field_3.sub' => [
                'sub error',
            ],
        ], $coll->toArray());
    }
}
