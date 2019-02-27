#shader vertex
#version 130

uniform vec3 u_Size;
uniform vec3 u_Position;

in vec3 vPos;

void main() {
    gl_Position = vec4(vPos * u_Size + u_Position, 1.0);
}

#shader fragment
#version 130

out vec4 fColor;

void main() {
    fColor = vec4(1.0, 1.0, 1.0, 1.0);
}
