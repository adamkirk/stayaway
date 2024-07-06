<?php

namespace App\Collections;

use Countable;
use ArrayAccess;
use Traversable;
use ArrayIterator;
use IteratorAggregate;
use App\Entities\Organisation;
use App\Collections\Collection;

/**
 * @method Organisation offsetGet()
 * @method Organisation[] all()
 * @property Organisation[] $items
 */
class OrganisationCollection extends Collection
{
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
}
