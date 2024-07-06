<?php

namespace App\ValueObjects;

use InvalidArgumentException;
use Ramsey\Uuid\UuidInterface;
use Ramsey\Uuid\Rfc4122\Fields;
use Ramsey\Uuid\Uuid as BaseUuid;
use Symfony\Component\Validator\Constraints\Uuid as UuidConstraint;

class Uuid
{
    const ASSERTION_TYPE = UuidConstraint::V7_MONOTONIC;

    protected function __construct(
        protected readonly UuidInterface $uuid,
    ) {}

    public function toString(): string
    {
        return $this->uuid->toString();
    }

    public function __toString(): string
    {
        return $this->toString();
    }

    public static function new()
    {
        return new self(BaseUuid::uuid7());
    }

    public static function fromString(string $input): self
    {
        $uuid = BaseUuid::fromString($input);

        $fields = new Fields($uuid->getFields()->getBytes());

        if ($fields->getVersion() !== 7) {
            throw new InvalidArgumentException("Must supply a uuid v7");
        }

        return new self(BaseUuid::fromString($uuid));
    }
}