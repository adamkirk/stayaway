<?php

namespace App\Collections;

use Countable;
use ArrayAccess;
use Traversable;
use ArrayIterator;
use IteratorAggregate;

/**
 * Assumes the class uses the items property for the array.
 */
abstract class Collection implements IteratorAggregate, ArrayAccess, Countable
{
    protected array $items;

    public static function fromArray(array $items): static
    {
        return new static(...$items);
    }

    public function all(): array
    {
        return $this->items;
    }

    public function isEmpty(): bool
    {
        return empty($this->items);
    }

    public function count(): int
    {
        return count($this->items);
    }

    public function getIterator(): Traversable
    {
        return new ArrayIterator($this->items);
    }

    public function offsetSet($offset, $value): void
    {
        if (is_null($offset)) {
            $this->items[] = $value;
        } else {
            $this->items[$offset] = $value;
        }
    }

    public function offsetExists($offset): bool
    {
        return isset($this->items[$offset]);
    }

    public function offsetUnset($offset): void
    {
        unset($this->items[$offset]);
    }

    public function offsetGet($offset): mixed
    {
        return isset($this->items[$offset]) ? $this->items[$offset] : null;
    }
}
