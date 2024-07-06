<?php

namespace App\Api\Responses;

use ReflectionClass;
use ReflectionProperty;
use App\Contracts\Arrayable;
use App\Api\Translation\TranslatesFieldNames;

trait ConvertsSelfToArray
{
    use TranslatesFieldNames;

    public function toArray(): array
    {
        $data = [];
        $reflClass = new ReflectionClass($this);
        foreach ($reflClass->getProperties(ReflectionProperty::IS_PUBLIC) as $prop) {
            $value = $this->{$prop->name};

            if ($value instanceof Arrayable) {
                $value = $value->toArray();
            }

            $data[$this->translate($prop->name)] = $value;
        }

        return $data;
    }
}