<?php

namespace App\Collections;

use Countable;
use ArrayAccess;
use IteratorAggregate;
use Illuminate\Support\Arr;
use App\Collections\Collection;
use App\Errors\ValidationError;
use App\Collections\ActsLikeArray;

/**
 * @method ValdiationError offsetGet()
 * @method ValdiationError[] all()
 * @property ValdiationError[] $items
 */
class ValidationErrorCollection extends Collection
{
    protected function __construct(
        ValidationError ...$items
    ) {
        $this->items = $items;
    }

    public function add(ValidationError $error)
    {
        $this->items[] = $error;
    }

    public static function new(ValidationError ...$errors)
    {
        return new self(...$errors);
    }

    public function toArray(): array
    {
        $translated = [];

        foreach ($this->items as $error) {
            $translated[$error->name()][] = $error->message();
        }

        return $translated;
    }
}