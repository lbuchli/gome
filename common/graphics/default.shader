#shader vertex
#version 130

uniform mat4 u_MVP;

in vec3 vertex_pos;
in vec2 vertex_uv;
in vec3 vertex_normal;

out vec3 normal;

void main() {
    gl_Position = u_MVP * vec4(vertex_pos, 1.0);
	normal = vertex_normal;
}

#shader fragment
#version 130

in vec3 normal;

out vec4 fColor;

void main() {
    fColor = vec4(normal + vec3(.5, .5, .5), 1.0);
}
