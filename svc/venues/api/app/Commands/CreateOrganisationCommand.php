<?php

namespace App\Commands;

use App\ValueObjects\Uuid;
use Illuminate\Http\Request;
use App\Entities\Organisation;
use App\Validation\Validatable;
use App\Validation\ValidatesByAttributes;
use Illuminate\Foundation\Bus\Dispatchable;
use App\Api\Requests\PopulatableFromRequest;
use App\Api\Translation\TranslatesFieldNames;
use App\Collections\ValidationErrorCollection;
use Symfony\Component\Validator\Constraints as Assert;
use Symfony\Component\Validator\Validator\ValidatorInterface;

class CreateOrganisationCommand implements PopulatableFromRequest, Validatable
{
    use TranslatesFieldNames;
    use ValidatesByAttributes;
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

    #[Assert\Regex(
        pattern: Organisation::SLUG_CHARACTER_SET,
        message: "The slug must start and end with a number or letter, and may contain letters, numbers and hyphens",
    )]
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

    public function validate(): ?ValidationErrorCollection
    {
        return $this->validateSelf();
    }

    public function populate(Request $request)
    {
        $this->name = $request->get($this->translate('name'), '');
        $this->slug = $request->get($this->translate('slug'), null);
    }

    protected function getValidator(): ValidatorInterface
    {
        return $this->validator;
    }
}