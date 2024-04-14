<?php

namespace App\Entities;

use App\ValueObjects\Uuid;
use Illuminate\Support\Str;
use App\ValueObjects\Organisation as VO;

class Organisation
{
    const EXPECTED_UUID_VERSION=7;

    protected function __construct(
        protected readonly Uuid $id,
        protected VO\Name $name,
        protected VO\Slug $slug,
    ) {}

    public function id(): Uuid
    {
        return $this->id;
    }

    public function setName(VO\Name $name): self
    {
        $this->name = $name;

        return $this;
    }
    
    public function name(): VO\Name
    {
        return $this->name;
    }

    public function setSlug(VO\Slug $slug): self
    {
        $this->slug = $slug;

        return $this;
    }

    public function slug(): VO\Slug
    {
        return $this->slug;
    }

    protected static function slugify(string $val)
    {
        return VO\Slug::new(Str::slug($val));
    }

    public static function new(
        Uuid $id,
        VO\Name $name,
        ?VO\Slug $slug
    ): self
    {
        if ($slug === null) {
            $slug = self::slugify($name->value());
        }
        
        return new self(
            id: $id,
            name: $name,
            slug: $slug,
        );
    }
}