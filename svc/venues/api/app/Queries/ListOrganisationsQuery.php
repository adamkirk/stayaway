<?php

namespace App\Queries;

use App\Queries\ListQuery;
use App\ValueObjects\Uuid;
use Illuminate\Http\Request;
use App\Queries\HasPagination;
use App\Queries\PaginationInput;
use App\Validation\Validatable;
use App\Api\Translation\HttpField;
use App\Validation\ValidatesByAttributes;
use App\Http\V1\Responses\NotFound;
use App\Api\Translation\FieldPlacement;
use App\Collections\ValidationErrorCollection;
use Illuminate\Foundation\Bus\Dispatchable;
use App\Api\Requests\PopulatableFromRequest;
use Illuminate\Contracts\Support\Responsable;
use App\Validation\ExposesPostValidationHook;
use App\Http\V1\Responses\BadRequestWithErrors;
use Symfony\Component\Validator\Constraints as Assert;
use Symfony\Component\Validator\Validator\ValidatorInterface;

class ListOrganisationsQuery implements PopulatableFromRequest, Validatable, ExposesPostValidationHook
{
    use ValidatesByAttributes;
    use Dispatchable;
    use HasPaginationAndOrdering;

    public function populate(Request $request)
    {
        $this->populatePaginationAndOrdering($request);
    }

    public function invalidResponse(ValidationErrorCollection $errors): Responsable
    {
        return BadRequestWithErrors::new($errors);
    }

    public function validate(ValidatorInterface $validator): ?ValidationErrorCollection
    {
        return $this->validateSelf($validator);
    }

    public function postValidationHook(): void
    {
        $this->paginationAndOrderingPostHook();
    }

    protected static function getDefaultOrderField(): string
    {
        return "name";
    }

    protected static function getDefaultOrderDirection(): string
    {
        return OrderDirection::ASC;
    }
}