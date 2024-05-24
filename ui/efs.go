package ui

import (
	"embed"
)

//go:embed "html" "static"
var Files embed.FS




// The important line here is the comment above. This lools like a comment, but it is actually a 'special directive'. When our app
// is compiled, this comment directive instructs Go to store the files from our ui/htlm and ui/static folders in an embed.fs embedded filesytem
// referenced by the gloabl variable Files.
// The comment directive must be placed immediately above the variable in which you want to store the embedded files.
// The directive has the general format go:em... <paths>, and it's OK to specify multiple paths in one directive. The paths should be relative to the 
// source code file containg the directive. So in our case, it embeds the directives ui/static and ui/html from our project.
// You can only use the go:em... directive on global varaibles at package level, not within functions or methods. If you try to do use it within
// a function or method, you'll get the error "go:em... cannot apply to var inside func" at compile time.
// Paths cannot contain "." or ".." elements, nor may they begin or end with a "/". This essentially restricts you to only embedding files that are
// contained in the same directory (or subdirectory) as the source code which has the go:em... directive.
// If a path is to a directory, then all files in that directory are recursively embedded, except for files with names that begin with "." or "_". 
// If you want to include these files you should use the "all:" prefix, like "go:em... "all:static".
// The path separator should always be a foward slash, even on Windows machines.
// The embedded file sytem is always rooted in the directory which contains the go:em... directive. So, in the example above, our Files varaible contains
// an 'embed.FS" embedded filesystem and the root of that filesystem is our ui directory.
