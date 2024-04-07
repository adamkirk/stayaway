<?php

namespace App\Repositories;

enum DeleteResult 
{
    case Deleted;
    case NotFound;
}