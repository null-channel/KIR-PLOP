<?php

use Illuminate\Support\Facades\Route;
use App\Http\Controllers\KubernetesController;

Route::get('/', function () {
    return view('welcome');
});

Route::get('/kubernetes', [KubernetesController::class, 'show'])->name('kubernetes.show');
Route::post('/kubernetes', [KubernetesController::class, 'submit'])->name('kubernetes.submit');
