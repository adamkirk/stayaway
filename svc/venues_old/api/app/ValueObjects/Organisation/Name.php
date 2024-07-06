<?php

namespace App\ValueObjects\Organisation;

use App\Exceptions\InvalidValueException;

class Name
{
    const MIN_LENGTH = 3;
    const MAX_LENGTH = 255;

    protected function __construct(
        protected string $name
    ) {}

    public static function guardValue(string $name): void
    {
        $length = strlen($name);
        if ($length < self::MIN_LENGTH) {
            throw new InvalidValueException('value is too short');
        }

        if ($length > self::MAX_LENGTH) {
            throw new InvalidValueException('value is too long');
        }
    }

    public function value(): string
    {
        return $this->name;
    }

    public static function new(string $name): self
    {
        self::guardValue($name);

        return new self($name);
    }
}