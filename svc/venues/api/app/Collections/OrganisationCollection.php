<?php

namespace App\Collections;

use App\Entities\Organisation;

class OrganisationCollection
{
    /** @var Organisation[] */
    protected array $items;

    protected function __construct(
        Organisation ...$items
    ) {
        $this->items = $items;
    }

    public static function new(Organisation ...$items): self
    {
        return new self(...$items);
    }

    public function add(Organisation $item): void
    {
        $this->items[] = $item;
    }

    /** @return Organisation[] */
    public function all(): array
    {
        return $this->items;
    }

    public static function fromArray(array $items): self
    {
        return new self(...$items);
    }
}