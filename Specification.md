# Language Specification

## Notation

The syntax is specified using a [variant](https://en.wikipedia.org/wiki/Wirth_syntax_notation) of Extended Backus-Naur Form (EBNF):

```ebnf
Syntax      = { Production } .
Production  = production_name "=" [ Expression ] "." .
Expression  = Term { "|" Term } .
Term        = Factor { Factor } .
Factor      = production_name | token [ "…" token ] | Group | Option | Repetition .
Group       = "(" Expression ")" .
Option      = "[" Expression "]" .
Repetition  = "{" Expression "}" .
```

Productions are expressions constructed from terms and the following operators, in increasing precedence:


```ebnf
|   alternation
()  grouping
[]  option (0 or 1 times)
{}  repetition (0 to n times)
```

The form a … b represents the set of characters from a through b as alternatives.

## Source Code Representation

Source code is Unicode text encoded in UTF-8. The text is not canonicalized, so a single accented code point is distinct from the same character constructed from combining an accent and a letter; those are treated as two code points. In this document, the term character will be used to refer to a Unicode code point in the source text.

Each code point is distinct; for instance, uppercase and lowercase letters are different characters.

Implementation restriction: A compiler *must* disallow the NUL character (U+0000) in the source text.
Implementation restriction: A compiler *may* ignore a UTF-8-encoded byte order mark (U+FEFF) if it is the first Unicode code point in the source text. A byte order mark *must* be disallowed anywhere else in the source text.

### Characters

The following terms are used to denote specific Unicode character categories:

```ebnf
newline        = /* the Unicode code point U+000A */ .
unicode_char   = /* an arbitrary Unicode code point except newline */ .
unicode_letter = /* a Unicode code point categorized as "Letter" */ .
unicode_digit  = /* a Unicode code point categorized as "Number, decimal digit" */ .
```

In [The Unicode Standard 8.0](https://www.unicode.org/versions/Unicode8.0.0/), Section 4.5 "General Category" defines a set of character categories. Hippo treats all characters in any of the Letter categories Lu, Ll, Lt, Lm, or Lo as Unicode letters, and those in the Number category Nd as Unicode digits.

### Letters and digits

The underscore character _ (U+005F) is considered a letter.

<!-- TODO: Add bindary, octal, hex etc. -->

```ebnf
letter        = unicode_letter | "_" .
decimal_digit = "0" … "9" .
```

## Lexical Elements

### Tokens

Tokens form the vocabulary of the Hippo language. There are four classes: *identifiers*, *keywords*, *operators* and *punctuation*, and *literals*. *White space*, formed from spaces (U+0020), horizontal tabs (U+0009), carriage returns (U+000D), and newlines (U+000A), is ignored except as it separates tokens that would otherwise combine into a single token. 
<!-- Also, a newline or end of file may trigger the insertion of a semicolon. While breaking the input into tokens, the next token is the longest sequence of characters that form a valid token. -->

### Identifiers

Identifiers name program entities such as variables and types. An identifier is a sequence of one or more letters and digits. The first character in an identifier must be a letter.

```ebnf
identifier = letter { letter | unicode_digit } .
```

```
a
_b9
αβ
```

Some identifiers are [predeclared](#predeclared).

### Keywords

The following keywords are reserved and may not be used as identifiers.

<!-- TODO: Add more keywords -->

```
fn var const
```

### Operators and punctuation

The following character sequences represent [operators](#operators) (including [assignment operators](#assignment-operators)) and punctuation:

```
+    &&    ==    !=    (     )
-    ||    <     <=    [     ]
*    <-    >     >=    {     }
/    =>    =     ...   ,     :
%    ++    !           .
```

### Literals

#### Integer literals

For readability, an underscore character _ may appear after a base prefix or between successive digits; such underscores do not change the literal's value.

<!-- TODO: Add floating point, binary etc -->

```ebnf
int_lit     = decimal_lit .
decimal_lit = decimal_digit { [ "_" ] decimal_digit } .
```

```
42
4_2
042
0_42
```

#### String literals

<!-- TODO: Add string literal spec -->

## Constants

A constant is a value denoted by an identifier. A constant may be a boolean, string, or numeric value. The boolean truth values are represented by the predeclared constants `true` and `false`. The predeclared identifier `null` denotes the zero value for a function.

## Variables

A variable is a storage location for holding a *value*. The set of permissible values is determined by the variable's [type](#types).

The *static type* of a variable is the type given at its declaration.
The *dynamic type* of a variable is the type of the value currently stored in the variable.

```
var x int = 1 // x is a variable of type int
var y = 2     // y is a variable of type int
var z int     // z is a variable of type int with the zero value for int
```

A variables value is retrieved by referring to the variable in an [expression](#expressions); it is set by assigning a value to the variable.

A variable declaration may be followed by an *initialization*, which specifies the variable is to be initialized to a particular value. If an initialization is present, the type can be omitted; the variable will take the type of the initializer.

If a variable declaration does not include an initialization, the variable is initialized with the [zero value](#zero-values) for its type.

## Types

TBC...