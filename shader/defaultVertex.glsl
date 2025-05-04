#version 330

precision highp float;

in  vec3 vertexPosition;
in  vec3 vertexNormal;

uniform mat4 mvp;
uniform mat4 model;
uniform mat4 matNormal;

out vec3 fragPos;
out vec3 fragNormal;

void main()
{
    fragPos    = vec3(model * vec4(vertexPosition, 1.0));

    mat3 normalMat = mat3(matNormal);
    // mat3 normalMat = transpose(inverse(mat3(model)));
    fragNormal = normalize(normalMat * vertexNormal);

    gl_Position = mvp * vec4(fragPos, 1.0);
}
