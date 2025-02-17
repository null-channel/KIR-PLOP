<!DOCTYPE html>
<html lang="{{ str_replace('_', '-', app()->getLocale()) }}">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Kubernetes Control</title>
    <!-- Fonts -->
    <link rel="preconnect" href="https://fonts.bunny.net">
    <link href="https://fonts.bunny.net/css?family=figtree:400,500,600&display=swap" rel="stylesheet" />
    <!-- Styles -->
    <style>
        body {
            font-family: 'Figtree', sans-serif;
            background-color: #f3f4f6;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            margin: 50px auto;
            padding: 20px;
        }
        .card {
            background-color: white;
            border-radius: 8px;
            padding: 20px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        h1 {
            color: #1a202c;
            text-align: center;
            margin-bottom: 30px;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            margin-bottom: 8px;
            color: #4a5568;
        }
        input[type="number"] {
            width: 100%;
            padding: 8px;
            border: 1px solid #e2e8f0;
            border-radius: 4px;
            font-size: 16px;
            box-sizing: border-box;
        }
        button {
            background-color: #4f46e5;
            color: white;
            padding: 10px 20px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            width: 100%;
            font-size: 16px;
        }
        button:hover {
            background-color: #4338ca;
        }
        .alert {
            padding: 10px;
            margin-bottom: 20px;
            border-radius: 4px;
        }
        .alert-success {
            background-color: #dcfce7;
            color: #166534;
        }
        .alert-error {
            background-color: #fee2e2;
            color: #991b1b;
        }
        form {
            width: 100%;
            margin: 0;
        }
        .form-group {
            margin-bottom: 20px;
            width: 100%;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="card">
            <h1>Kubernetes Control Panel</h1>
            
            @if(session('success'))
                <div class="alert alert-success">
                    {{ session('success') }}
                </div>
            @endif

            @if(session('error'))
                <div class="alert alert-error">
                    {{ session('error') }}
                </div>
            @endif

            <form action="{{ route('kubernetes.submit') }}" method="POST">
                @csrf
                <div class="form-group">
                    <label for="name">Resource Name:</label>
                    <input type="text" id="name" name="name" required value="{{ old('name') }}" pattern="[a-z0-9]([-a-z0-9]*[a-z0-9])?" title="Name must consist of lowercase alphanumeric characters or '-', and must start and end with an alphanumeric character">
                </div>
                <div class="form-group">
                    <label for="number">Enter Number of vertices:</label>
                    <input type="number" id="number" name="number" required min="1" value="{{ old('number') }}">
                </div>
                <button type="submit">Send to Kubernetes</button>
            </form>
        </div>
    </div>
</body>
</html> 