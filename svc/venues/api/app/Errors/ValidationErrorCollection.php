<?php

namespace App\Errors;

use Illuminate\Support\Arr;

class ValidationErrorCollection
{
    /** @var ValidationError[] $errors */
    protected array $errors;

    protected function __construct(
        ValidationError ...$errors
    ) {
        $this->errors = $errors;
    }

    public function add(ValidationError $error)
    {
        $this->errors[] = $error;
    }

    public function isEmpty(): bool
    {
        return empty($this->errors);
    }

    /** @return ValidationError[] */
    public function all(): array
    {
        return $this->errors;
    }

    public static function new(array $errors = [])
    {
        return new self(...$errors);
    }

    public function toArray(): array
    {
        $translated = [];

        foreach ($this->errors as $error) {
            $translated[$error->name()][] = $error->message();
        }

        return $translated;
    }
}