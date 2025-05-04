#version 330

precision highp float;

in  vec3 vertexPosition;

uniform mat4 mvp;
uniform mat4 model;
uniform vec3 lightPos;

out vec3 fragPos;
flat out int hideFrag;

void main()
{
    fragPos    = vec3(model * vec4(vertexPosition, 1.0));
    vec3 lightDir = normalize( fragPos - lightPos);

    fragPos = fragPos + lightDir * (-fragPos.y / lightDir.y) + .00001;
    if (fragPos.y < 0.0) {
        hideFrag = 1;
    }
    else {
        hideFrag = 0;
    }

    gl_Position = mvp * vec4(fragPos, 1.0);
}
