package main

var HTML string = `
<!DOCTYPE html>

<html lang="en">

<head>
    <meta charset="utf-8">

	<style>
html {
    background-color: #282828;
}

body {
    margin: 0 auto;
    max-width: 23cm;
    padding: 1em;
    background-color: #fdfdfd;
}

nav {
	padding-left: 3em;
}

h1 {
    font-size: 2.9em;
    font-family: Courier, monospace;
    color: #af111c;
}

h2 {
    font-size: 1.8em;
}

ul {
    list-style-type: none;
}

.item-days {
    font-family: Helvetica, sans-serif;
}

.feed-title {
    font-weight: bold;
}

.item-entries {
    font-family: Palatino, Georgia, serif;
    font-size: 18px;
    line-height: 1.5;
}

.list-entries {
	max-width: 85%;
}
	</style>

    <title>Evil feed reader</title>
</head>

<body>
    <nav>
        <h1>Evil feed reader</h1>
    </nav>

    <main>
        <ul class="list-days">
        {{ range .Days }}
        <li class="item-days">
            <h2>{{ .Date.Day }}&nbsp;{{ .Date.Month }}&nbsp;{{ .Date.Year }}</h2>
        </li>
        <ul class="list-entries">
            {{ range .Entries }}
            <li class="item-entries">
                <p>
                <span class="feed-title">{{ .Feed }}</span>:&ensp;
                <span class="item-title">
                    <a href="{{ .Link }}">{{ .Title }}</a>
                </span>
                </p>
            </li>
            {{ end  }}
        </ul>
        {{ end  }}
        </ul>
    </main>

</body>

</html>
`
