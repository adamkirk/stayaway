<?php

namespace App\ValueObjects\Organisation;

use App\Exceptions\InvalidValueException;

class Slug
{
    const MIN_LENGTH = 2;
    const MAX_LENGTH = 255;
    const CHARACTER_SET = '/^[a-z0-9]{1}([a-z0-9\-])*[a-z0-9]{1}$/';

    protected function __construct(
        protected string $slug
    ) {}

    public static function guardValue(string $slug): void
    {
        $length = strlen($slug);

        if ($length < self::MIN_LENGTH) {
            throw new InvalidValueException('value is too short');
        }

        if ($length > self::MAX_LENGTH) {
            throw new InvalidValueException('value is too long');
        }

        if (preg_match(self::CHARACTER_SET, $slug) !== 1) {
            throw new InvalidValueException('value contains invalid characters');
        }
    }

    public function value(): string
    {
        return $this->slug;
    }

    public static function new(string $slug): self
    {
        self::guardValue($slug);

        return new self($slug);
    }
}