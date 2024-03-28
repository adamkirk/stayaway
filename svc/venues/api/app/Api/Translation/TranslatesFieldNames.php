<?php

namespace App\Api\Translation;

use ReflectionProperty;
use Illuminate\Support\Str;
use App\Api\Translation\HttpField;

trait TranslatesFieldNames
{
    protected function translate(string $propName, object|string|null $on = null): string
    {
        $prop = new ReflectionProperty($on ?? $this, $propName);
        $attr = $prop->getAttributes(HttpField::class)[0] ?? null;

        if ($attr !== null) {
            $instance = $attr->newInstance();

            if ($instance->name() !== null) {
                return $instance->name();
            }
        }

        return Str::snake($propName);
    }
}