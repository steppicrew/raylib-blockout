#version 330

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
    fragNormal = (matNormal * vec4(vertexNormal, 0.0)).xyz;

    mat3 normalMat = mat3(matNormal);
    fragNormal = normalize(normalMat * vertexNormal);

    gl_Position = mvp * vec4(vertexPosition, 1.0);
}
