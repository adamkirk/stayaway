<?php

use Illuminate\Support\Facades\Route;
use App\Http\V1\Controllers\OrganisationsController;

Route::post('/v1/organisations', [OrganisationsController::class, 'create']);
Route::get('/v1/organisations', [OrganisationsController::class, 'list']);
Route::get('/v1/organisations/{organisation_id}', [OrganisationsController::class, 'get']);
Route::delete('/v1/organisations/{organisation_id}', [OrganisationsController::class, 'delete']);