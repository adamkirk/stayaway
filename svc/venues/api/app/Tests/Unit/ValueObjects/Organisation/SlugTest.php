<?php

namespace App\Tests\Unit\ValueObjects\Organisation;

use Tests\UnitTestCase;
use App\ValueObjects\Organisation\Slug;
use App\Exceptions\InvalidValueException;
use PHPUnit\Framework\Attributes\DataProvider;

class SlugTest extends UnitTestCase
{
    #[DataProvider('constructionFailures')]
    public function test_that_construction_only_allows_valid_properties(string $slug, string $expectedFailureReason): void
    {
        $this->expectExceptionObject(new InvalidValueException(
            $expectedFailureReason,
        ));

        Slug::new($slug);
    }

    public static function constructionFailures(): array
    {
        return [
            'slug_too_short' => [
                str_repeat("X", Slug::MIN_LENGTH - 1),
                'value is too short',
            ],
            'slug_too_long' => [
                str_repeat("X", Slug::MAX_LENGTH + 1),
                'value is too long',
            ],
            'slug_invalid_characters' => [
                'I cant have spaces brah!',
                'value contains invalid characters',
            ],
            'slug_invalid_characters_specials' => [
                'cheeky-$-sign',
                'value contains invalid characters',
            ],
            'slug_invalid_characters_underscore' => [
                'underscore_billy',
                'value contains invalid characters',
            ],
        ];
    }
}