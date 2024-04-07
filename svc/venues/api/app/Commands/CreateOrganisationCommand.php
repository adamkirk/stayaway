<?php

namespace App\Commands;

use App\ValueObjects\Uuid;
use Illuminate\Http\Request;
use App\Entities\Organisation;
use App\Api\Requests\Validatable;
use App\Api\Requests\ValidatesSelf;
use Illuminate\Foundation\Bus\Dispatchable;
use App\Api\Requests\PopulatableFromRequest;
use App\Api\Translation\TranslatesFieldNames;
use Symfony\Component\Validator\Constraints as Assert;
use Symfony\Component\Validator\Validator\ValidatorInterface;

class CreateOrganisationCommand implements PopulatableFromRequest, Validatable
{
    use TranslatesFieldNames;
    use ValidatesSelf {
        validate as protected validateSelf;
    }
    use Dispatchable;

    public readonly Uuid $generatedId;

    #[Assert\NotBlank]
    #[Assert\Length(
        min: Organisation::NAME_MIN_LENGTH,
        max: Organisation::NAME_MAX_LENGTH,
        minMessage: 'The name must be at least {{ limit }} characters long',
        maxMessage: 'The name cannot be longer than {{ limit }} characters',
    )]
    public readonly string $name;

    #[Assert\Regex(Organisation::SLUG_CHARACTER_SET)]
    #[Assert\Length(
        min: Organisation::SLUG_MIN_LENGTH,
        max: Organisation::SLUG_MAX_LENGTH,
        minMessage: 'The slug must be at least {{ limit }} characters long',
        maxMessage: 'The slug cannot be longer than {{ limit }} characters',
    )]
    public readonly ?string $slug;
    
    public function __construct(
        protected ValidatorInterface $validator
    ) {
        $this->generatedId = Uuid::new();
    }

    public function populate(Request $request)
    {
        $this->name = $request->get($this->translate('name'), null);
        $this->slug = $request->get($this->translate('slug'), null);
    }

    protected function getValidator(): ValidatorInterface
    {
        return $this->validator;
    }
}