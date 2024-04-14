<?php

namespace App\Commands;

use App\ValueObjects\Uuid;
use Illuminate\Http\Request;
use App\Buses\DefinesHandler;
use App\Entities\Organisation;
use App\Validation\Validatable;
use App\ValueObjects\Organisation as VO;
use App\Validation\ValidatesByAttributes;
use App\Handlers\CreateOrganisationHandler;
use Illuminate\Foundation\Bus\Dispatchable;
use App\Api\Requests\PopulatableFromRequest;
use App\Api\Translation\TranslatesFieldNames;
use App\Collections\ValidationErrorCollection;
use Symfony\Component\Validator\Constraints as Assert;
use Symfony\Component\Validator\Validator\ValidatorInterface;

class CreateOrganisationCommand implements PopulatableFromRequest, Validatable, DefinesHandler
{
    use TranslatesFieldNames;
    use ValidatesByAttributes;
    use Dispatchable;

    protected readonly Uuid $generatedId;

    #[Assert\NotBlank]
    #[Assert\Length(
        min: VO\Name::MIN_LENGTH,
        max: VO\Name::MAX_LENGTH,
        minMessage: 'The name must be at least {{ limit }} characters long',
        maxMessage: 'The name cannot be longer than {{ limit }} characters',
    )]
    protected readonly string $name;

    #[Assert\Regex(
        pattern: VO\Slug::CHARACTER_SET,
        message: "The slug must start and end with a number or letter, and may contain letters, numbers and hyphens",
    )]
    #[Assert\Length(
        min: VO\Slug::MIN_LENGTH,
        max: VO\Slug::MAX_LENGTH,
        minMessage: 'The slug must be at least {{ limit }} characters long',
        maxMessage: 'The slug cannot be longer than {{ limit }} characters',
    )]
    protected readonly ?string $slug;
    
    public function __construct(
        protected ValidatorInterface $validator
    ) {
        $this->generatedId = Uuid::new();
    }

    public static function getHandler(): string
    {
        return CreateOrganisationHandler::class;
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

    public function id(): Uuid
    {
        return $this->generatedId;
    }

    public function name(): VO\Name
    {
        return VO\Name::new($this->name);
    }

    public function slug(): ?VO\Slug
    {
        return $this->slug !== null ? VO\Slug::new($this->slug) : null;
    }

    protected function getValidator(): ValidatorInterface
    {
        return $this->validator;
    }
}