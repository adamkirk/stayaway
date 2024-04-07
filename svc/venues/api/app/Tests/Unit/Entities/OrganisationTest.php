<?php

namespace App\Tests\Unit\Entities;

use Tests\UnitTestCase;
use App\ValueObjects\Uuid;
use App\Entities\Organisation;
use App\Exceptions\InvalidPropertyException;
use PHPUnit\Framework\Attributes\DataProvider;

class OrganisationTest extends UnitTestCase
{
    protected function mockUuid(): Uuid
    {
        return $this->getMockBuilder(Uuid::class)
            ->disableOriginalConstructor()
            ->getMock();
    }

    #[DataProvider('constructionFailures')]
    public function test_that_construction_only_allows_valid_properties(string $name, string $slug, string $expectedPropertyFailure, string $expectedFailureReason): void
    {
        $this->expectExceptionObject(new InvalidPropertyException(
            Organisation::class,
            $expectedPropertyFailure,
            $expectedFailureReason,
        ));
        Organisation::new($this->mockUuid(), $name, $slug);
    }

    public static function constructionFailures(): array
    {
        return [
            'name_too_short' => [
                str_repeat("X", Organisation::NAME_MIN_LENGTH - 1),
                '',
                'name',
                'too short',
            ],
            'slug_too_short' => [
                'Valid Name',
                str_repeat("X", Organisation::SLUG_MIN_LENGTH - 1),
                'slug',
                'too short',
            ],
            'name_too_long' => [
                // Purposely didn't use the constant here, tests should be dumb i think and not auto-fix when code changes
                // The db will break if it gets bigger so we need to be sure.
                str_repeat("X", Organisation::NAME_MAX_LENGTH + 1),
                'valid-slug',
                'name',
                'too long',
            ],
            'slug_too_long' => [
                'Valid Name',
                // Purposely didn't use the constant here, tests should be dumb i think and not auto-fix when code changes
                // The db will break if it gets bigger so we need to be sure.
                str_repeat("X", Organisation::SLUG_MAX_LENGTH + 1),
                'slug',
                'too long',
            ],
            'slug_invalid_characters' => [
                'Valid Name',
                // Purposely didn't use the constant here, tests should be dumb i think and not auto-fix when code changes
                // The db will break if it gets bigger so we need to be sure.
                'I cant have spaces brah!',
                'slug',
                'invalid characters',
            ],
            'slug_invalid_characters_specials' => [
                'Valid Name',
                // Purposely didn't use the constant here, tests should be dumb i think and not auto-fix when code changes
                // The db will break if it gets bigger so we need to be sure.
                'cheeky-$-sign',
                'slug',
                'invalid characters',
            ],
            'slug_invalid_characters_underscore' => [
                'Valid Name',
                // Purposely didn't use the constant here, tests should be dumb i think and not auto-fix when code changes
                // The db will break if it gets bigger so we need to be sure.
                'underscore_billy',
                'slug',
                'invalid characters',
            ],
        ];
    }

    public function test_successful_construction(): void
    {
        $uuid = $this->mockUuid();
        $name = "A Valid Name";
        $slug = "a-valid-slug";
        $org = Organisation::new($uuid, $name, $slug);

        $this->assertInstanceOf(Organisation::class, $org);
    }


    public function test_slug_defaults(): void
    {
        $uuid = $this->mockUuid();
        $name = "A Valid Name !_special_chars";
        $org = Organisation::new($uuid, $name, null);

        $this->assertInstanceOf(Organisation::class, $org);
        $this->assertEquals("a-valid-name-special-chars", $org->slug());
    }


    public function test_getters_and_setters(): void
    {
        $uuid = $this->mockUuid();
        $name = "A Valid Name";
        $slug = "a-valid-slug";
        $org = Organisation::new($uuid, $name, $slug);

        $this->assertInstanceOf(Organisation::class, $org);
        $this->assertSame($uuid, $org->id());
        $this->assertEquals($name, $org->name());
        $this->assertEquals($slug, $org->slug());

        $newName = "A New Name";
        $org->setName($newName);
        $this->assertEquals($newName, $org->name());

        $newSlug = "a-new-slug";
        $org->setSlug($newSlug);
        $this->assertEquals($newSlug, $org->slug());
    }

    #[DataProvider('invalidNames')]
    public function test_name_setter_is_guarded(string $newName, string $expectedFailureReason): void
    {
        
        $org = Organisation::new($this->mockUuid(), "Valid Name", "valid-slug");
        $this->assertInstanceOf(Organisation::class, $org);
        $this->expectExceptionObject(new InvalidPropertyException(
            Organisation::class,
            'name',
            $expectedFailureReason,
        ));
        
        $org->setName($newName);
    }

    public static function invalidNames(): array
    {
        return [
            'empty' => [
                '',
                'too short',
            ],
            'too_short' => [
                str_repeat("X", Organisation::NAME_MIN_LENGTH - 1),
                'too short',
            ],
            'too_long' => [
                str_repeat("X", Organisation::NAME_MAX_LENGTH + 1),
                'too long',
            ],
        ];
    }

    #[DataProvider('invalidSlugs')]
    public function test_slug_setter_is_guarded(string $newSlug, string $expectedFailureReason): void
    {
        
        $org = Organisation::new($this->mockUuid(), "Valid Name", "valid-slug");
        $this->assertInstanceOf(Organisation::class, $org);
        $this->expectExceptionObject(new InvalidPropertyException(
            Organisation::class,
            'slug',
            $expectedFailureReason,
        ));
        
        $org->setSlug($newSlug);
    }

    public static function invalidSlugs(): array
    {
        return [
            'empty' => [
                '',
                'too short',
            ],
            'too_short' => [
                str_repeat("x", Organisation::SLUG_MIN_LENGTH - 1),
                'too short',
            ],
            'too_long' => [
                str_repeat("x", Organisation::SLUG_MAX_LENGTH + 1),
                'too long',
            ],
            'invalid_characters_underscore' => [
                'with_underscore',
                'invalid characters',
            ],
            'with_space' => [
                'with space',
                'invalid characters',
            ],
            'with_ampersand' => [
                'with-special-&',
                'invalid characters',
            ],
            'with_slash' => [
                'with-/-special',
                'invalid characters',
            ],
            'with_question_mark' => [
                'with-?-special',
                'invalid characters',
            ],
        ];
    }


}
