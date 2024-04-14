<?php

namespace App\Tests\Unit\ValueObjects\Organisation;

use Tests\UnitTestCase;
use App\ValueObjects\Organisation\Name;
use App\Exceptions\InvalidValueException;
use PHPUnit\Framework\Attributes\DataProvider;

class NameTest extends UnitTestCase
{
    #[DataProvider('constructionFailures')]
    public function test_that_construction_only_allows_valid_properties(string $name, string $expectedFailureReason): void
    {
        $this->expectExceptionObject(new InvalidValueException(
            $expectedFailureReason,
        ));

        Name::new($name);
    }

    public static function constructionFailures(): array
    {
        return [
            'name_too_short' => [
                str_repeat("X", Name::MIN_LENGTH - 1),
                'value is too short',
            ],
            'name_too_long' => [
                str_repeat("X", Name::MAX_LENGTH + 1),
                'value is too long',
            ],
        ];
    }
}