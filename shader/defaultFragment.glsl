#version 330

precision highp float;

in vec3 fragPos;
in vec3 fragNormal;

uniform vec3 viewPos;
uniform vec3 lightPos;

uniform vec3 lightColor;
uniform vec4 objectColor;

out vec4 finalColor;

void main()
{
    // Ambient
    float ambientStrength = 0.2;
    vec3 ambient = ambientStrength * lightColor;

    // Diffuse
    vec3 norm = normalize(fragNormal);
    vec3 lightDir = normalize(lightPos - fragPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * lightColor;

    // Specular
    float specularStrength = 4.0;
    vec3 viewDir = normalize(viewPos - fragPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), 32);
    vec3 specular = specularStrength * spec * lightColor;

    finalColor = vec4((ambient + diffuse + specular) * vec3(objectColor), objectColor.a);
    // finalColor = (ambient + diffuse) * objectColor;
    // finalColor = (ambient + specular) * objectColor;
    // finalColor = vec4(normalize(fragPos) * 0.5 + 0.5, 1.0);
}
