/*
From https://github.com/gojp/kana
Copyright 2013 Herman Schaaf and Shawn Smith

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package passgen

// HiraganaTable maps romaji to hiragana
var HiraganaTable = `	a	i	u	e	o	n
	あ	い	う	え	お	ん
x	ぁ	ぃ	ぅ	ぇ	ぉ
k	か	き	く	け	こ
ky	きゃ		きゅ		きょ
s	さ		す	せ	そ
sh	しゃ	し	しゅ		しょ
t	た		つ	て	と
ts			つ
ch	ちゃ	ち	ちゅ	ちぇ	ちょ
n	な	に	ぬ	ね	の
ny	にゃ		にゅ		にょ
h	は	ひ	ふ	へ	ほ
hy	ひゃ		ひゅ		ひょ
f			ふ
m	ま	み	む	め	も
my	みゃ		みゅ		みょ
y	や		ゆ		よ
r	ら	り	る	れ	ろ
ry	りゃ	りぃ	りゅ	りぇ	りょ
w	わ	ゐ		ゑ	を
g	が	ぎ	ぐ	げ	ご
gy	ぎゃ		ぎゅ		ぎょ
z	ざ		ず	ぜ	ぞ/ぢょ
j	じゃ/ぢゃ	じ/ぢ	じゅ/ぢゅ		じょ
d	だ		づ	で	ど
b	ば	び	ぶ	べ	ぼ
by	びゃ		びゅ		びょ
p	ぱ	ぴ	ぷ	ぺ	ぽ
py	ぴゃ		ぴゅ		ぴょ
v			ゔ`
