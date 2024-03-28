<?php

namespace App\Entities;

use App\ValueObjects\Uuid;
use Illuminate\Support\Str;
use App\Exceptions\InvalidPropertyException;

class Organisation
{
    const EXPECTED_UUID_VERSION=7;
    const NAME_MIN_LENGTH = 3;
    const SLUG_MIN_LENGTH = 3;
    const NAME_MAX_LENGTH = 255;
    const SLUG_MAX_LENGTH = 255;
    const SLUG_CHARACTER_SET = '/^[A-Za-z0-9\-]+$/';

    protected function __construct(
        protected readonly Uuid $id,
        protected string $name,
        protected string $slug,
    ) {}

    public function id(): Uuid
    {
        return $this->id;
    }

    public function setName(string $name): self
    {
        self::guardName($name);

        $this->name = $name;

        return $this;
    }

    public static function guardName(string $name): void
    {
        $length = strlen($name);
        if ($length < self::NAME_MIN_LENGTH) {
            throw new InvalidPropertyException(self::class, 'name', 'too short');
        }

        if ($length > self::NAME_MAX_LENGTH) {
            throw new InvalidPropertyException(self::class, 'name', 'too long');
        }
    }
    
    public function name(): string
    {
        return $this->name;
    }

    public function setSlug(string $slug): self
    {
        self::guardSlug($slug);

        $this->slug = $slug;

        return $this;
    }

    protected static function guardSlug(string $slug): void
    {
        $length = strlen($slug);

        if ($length < self::SLUG_MIN_LENGTH) {
            throw new InvalidPropertyException(self::class, 'slug', 'too short');
        }

        if ($length > self::SLUG_MAX_LENGTH) {
            throw new InvalidPropertyException(self::class, 'slug', 'too long');
        }

        if (preg_match(self::SLUG_CHARACTER_SET, $slug) !== 1) {
            throw new InvalidPropertyException(self::class, 'slug', 'invalid characters');
        }
    }

    public function slug(): string
    {
        return $this->slug;
    }

    protected static function slugify(string $val)
    {
        return Str::slug($val);
    }

    public static function new(
        Uuid $id,
        string $name,
        ?string $slug
    ): self
    {
        self::guardName($name);

        if ($slug === null) {
            $slug = self::slugify($name);
        }
        
        self::guardSlug($slug);

        return new self(
            id: $id,
            name: $name,
            slug: $slug,
        );
    }
}