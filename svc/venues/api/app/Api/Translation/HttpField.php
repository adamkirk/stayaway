<?php

namespace App\Api\Translation;

use Attribute;
use App\Api\Translation\FieldPlacement;

#[Attribute(Attribute::TARGET_PROPERTY)]
class HttpField
{
    public function __construct(
        protected string|null $name = null,
        protected FieldPlacement $in = FieldPlacement::Body,
    ) {}

    public function setName(string $name): void
    {
        $this->name = $name;
    }

    public function name(): string|null
    {
        return $this->name;
    }

    public function setPlacement(FieldPlacement $placement): void
    {
        $this->in = $placement;
    }

    public function placement(): FieldPlacement
    {
        return $this->in;
    }
}