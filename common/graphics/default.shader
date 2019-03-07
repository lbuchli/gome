#shader vertex
#version 130

uniform mat4 u_MVP;

in vec3 vertex_pos;
in vec2 vertex_uv;
in vec3 vertex_normal;

out vec2 uv;
out vec3 normal;

void main() {
	uv = vertex_uv;
	normal = vertex_normal;
    gl_Position = u_MVP * vec4(vertex_pos, 1.0);
}

#shader fragment
#version 130

in vec2 uv;
in vec3 normal;

out vec4 fColor;

void main() {
    fColor = vec4(1.0, normal * .5 + .5);
}
