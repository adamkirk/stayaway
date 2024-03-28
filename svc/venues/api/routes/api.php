<?php

use Illuminate\Support\Facades\Route;
use App\Http\V1\Controllers\OrganisationsController;


Route::post('/v1/organisations', [OrganisationsController::class, 'create']);