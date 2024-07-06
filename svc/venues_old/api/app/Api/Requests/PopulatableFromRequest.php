<?php

namespace App\Api\Requests;

use Illuminate\Http\Request;

interface PopulatableFromRequest
{
    public function populate(Request $request);
}