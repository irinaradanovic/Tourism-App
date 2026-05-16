package com.tourism.stakeholders;

import com.tourism.stakeholders.model.User;
import com.tourism.stakeholders.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.boot.CommandLineRunner;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class DataInitializer implements CommandLineRunner {

    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;

    @Override
    public void run(String... args) {
        if (userRepository.count() == 0) {
            
            User admin = new User();
            admin.setUsername("admin");
            admin.setEmail("admin@tourism.com");
            admin.setPassword(passwordEncoder.encode("admin123"));
            admin.setRole(User.Role.ADMIN);
            admin.setBlocked(false);
            userRepository.save(admin);
            System.out.println(">>> Admin kreiran: username=admin, password=admin123 (ID: 1)");

            String encodedPassword = passwordEncoder.encode("password");
            
            String[] firstNames = {"Marko", "Jelena", "Nikola", "Milica", "Stefan", "Ana", "Luka", "Jovana", "Igor", "Sara", "Filip", "Teodora", "Pavle", "Tara", "Nemanja"};
            String[] lastNames = {"Markovic", "Popovic", "Nikolic", "Petrovic", "Stankovic", "Jovanovic", "Lukic", "Djordjevic", "Ilic", "Kostic", "Mitrovic", "Antic", "Lazarevic", "Tasic", "Ristic"};

            for (int i = 0; i < 15; i++) {
                User user = new User();
                String username = firstNames[i].toLowerCase() + (i + 2);
                
                user.setUsername(username);
                user.setEmail(username + "@tourism.com");
                user.setPassword(encodedPassword); 
                user.setFirstName(firstNames[i]);
                user.setLastName(lastNames[i]);
                user.setRole(i % 2 == 0 ? User.Role.TOURIST : User.Role.GUIDE);
                user.setBlocked(false);
                user.setBiography("Ja sam " + firstNames[i]);
                user.setMotto("Moto " + firstNames[i]);
                user.setProfileImage(null);

                userRepository.save(user);
                System.out.println(">>> Kreiran korisnik: @" + username + " (Role: " + user.getRole() + ", ID: " + (i + 2) + ")");
            }
            System.out.println(">>> Uspešno inicijalizovano svih 16 korisnika (1 Admin + 15 Test korisnika)!");
        }
    }
}