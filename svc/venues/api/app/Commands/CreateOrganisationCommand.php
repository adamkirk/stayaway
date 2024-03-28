<?php

namespace App\Commands;

use App\ValueObjects\Uuid;
use Illuminate\Http\Request;
use App\Api\Requests\Validatable;
use App\Api\Requests\ValidatesSelf;
use Illuminate\Foundation\Bus\Dispatchable;
use App\Api\Requests\PopulatableFromRequest;

use App\Api\Translation\TranslatesFieldNames;
use Symfony\Component\Validator\Constraints as Assert;
use Symfony\Component\Validator\Validator\ValidatorInterface;


// Left the validation limits as hard-coded here so to not couple this to constants
// in the domain which may change. The API spec shouldn't necessarily 
// change with the domain rules, this is why it's versioned.
class CreateOrganisationCommand implements PopulatableFromRequest, Validatable
{
    use TranslatesFieldNames;
    use ValidatesSelf;
    use Dispatchable;

    public readonly Uuid $generatedId;

    #[Assert\NotBlank]
    #[Assert\Length(
        min: 3,
        max: 255,
        minMessage: 'The name must be at least {{ limit }} characters long',
        maxMessage: 'The name cannot be longer than {{ limit }} characters',
    )]
    public readonly string $name;

    #[Assert\Regex('/^[A-Za-z0-9\-]+$/')]
    #[Assert\Length(
        min: 3,
        max: 255,
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
        $this->name = $request->get($this->translate('name'), '');
        $this->slug = $request->get($this->translate('slug'), '');
    }

    protected function getValidator(): ValidatorInterface
    {
        return $this->validator;
    }
}