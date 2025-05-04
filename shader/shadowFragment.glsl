#version 330

precision highp float;

uniform vec4 objectColor;

flat in int hideFrag;

out vec4 finalColor;

void main()
{
    if (hideFrag == 1) {
        discard;
    }
    else {
        finalColor= objectColor;
    }
}
